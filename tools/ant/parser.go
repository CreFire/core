package antnet

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/vmihailenco/msgpack"
	"reflect"
)

type IMsgParser interface {
	C2S() interface{}
	S2C() interface{}
	C2SData() []byte
	S2CData() []byte
	C2SString() string
	S2CString() string
	setC2S(interface{})
}

type MsgParser struct {
	s2c     interface{}
	c2s     interface{}
	c2sFunc ParseFunc
	s2cFunc ParseFunc
	parser  IParser
}

func (r *MsgParser) setC2S(c2s interface{}) {
	r.c2s = c2s
}

func (r *MsgParser) C2S() interface{} {
	if r.c2s == nil && r.c2sFunc != nil {
		r.c2s = r.c2sFunc()
	}
	return r.c2s
}

func (r *MsgParser) S2C() interface{} {
	if r.s2c == nil && r.s2cFunc != nil {
		r.s2c = r.s2cFunc()
	}
	return r.s2c
}

func (r *MsgParser) C2SData() []byte {
	return r.parser.PackMsg(r.C2S())
}

func (r *MsgParser) S2CData() []byte {
	return r.parser.PackMsg(r.S2C())
}

func (r *MsgParser) C2SString() string {
	return string(r.C2SData())
}

func (r *MsgParser) S2CString() string {
	return string(r.S2CData())
}

type ParserType int

const (
	ParserTypePB  ParserType = iota //protobuf类型，用于和客户端交互
	ParserTypeCmd                   //cmd类型，类似telnet指令，用于直接和程序交互
	ParserTypeRaw                   //不做任何解析
)

type ParseErrType int

const (
	ParseErrTypeSendRemind ParseErrType = iota //消息解析失败发送提醒消息
	ParseErrTypeContinue                       //消息解析失败则跳过本条消息
	ParseErrTypeAlways                         //消息解析失败依然处理
	ParseErrTypeClose                          //消息解析失败则关闭连接
)

type ParseFunc func() interface{}

type IParser interface {
	GetType() ParserType
	GetErrType() ParseErrType
	ParseC2S(msg *Message) (IMsgParser, error)
	PackMsg(v interface{}) []byte
}

type MessageID struct {
	Cmd uint8 // cmd
	Act uint8 // action
}

func (m *MessageID) ID() int {
	return CmdAct(m.Cmd, m.Act)
}

type Parser struct {
	Type    ParserType
	ErrType ParseErrType

	msgMap  map[int]MsgParser
	cmdRoot *cmdParseNode
	parser  IParser
	typeMap map[reflect.Type]*MessageID //Key type msgId
}

func (r *Parser) NewMessage(msg proto.Message) *Message {
	msgType := reflect.TypeOf(msg).Elem()
	msgId, ok := r.typeMap[msgType]
	if !ok {
		panic(Sprintf("PbParser NewSendMessage typeMap failed: %v", msgType))
	}
	parser, ok := r.msgMap[msgId.ID()]
	if !ok {
		panic(Sprintf("PbParser NewSendMessage msgMap failed: %v", msgType))
	}
	return NewPbMsg(msgId.Cmd, msgId.Act, msg, &parser)
}

func (r *Parser) Deserialize(cmd uint8, act uint8, data []byte) (*Message, error) {
	msg := NewMsg(uint8(cmd), uint8(act), data)
	p, ok := r.msgMap[CmdAct(cmd, act)]
	if ok {
		if p.C2S() != nil {
			err := PBUnPack(data, p.C2S())
			if err != nil {
				return nil, NewErr("Deserialize PBUnPack error cmd: %d, act: %d, err: %v", msg.Head.Cmd, msg.Head.Act, err)
			}
		} else {
			return nil, NewErr("Deserialize C2S is nil: cmd: %d, act: %d", cmd, act)
		}
	} else {
		return nil, NewErr("Deserialize not found cmd: %d, act: %d", cmd, act)
	}
	msg.IMsgParser = &p
	return msg, nil
}

func (r *Parser) Get() IParser {
	switch r.Type {
	case ParserTypePB:
		if r.parser == nil {
			r.parser = &pBParser{Parser: r}
		}
	case ParserTypeCmd:
		return &cmdParser{Parser: r}
	case ParserTypeRaw:
		return nil
	}

	return r.parser
}

func (r *Parser) GetType() ParserType {
	return r.Type
}

func (r *Parser) GetErrType() ParseErrType {
	return r.ErrType
}

func (r *Parser) RegisterFunc(cmd, act uint8, c2sFunc ParseFunc, s2cFunc ParseFunc) {
	if r.msgMap == nil {
		r.msgMap = map[int]MsgParser{}
	}
	r.msgMap[CmdAct(cmd, act)] = MsgParser{c2sFunc: c2sFunc, s2cFunc: s2cFunc}
}

func (r *Parser) Register(cmd uint8, act uint8, c2s interface{}, s2c interface{}) {
	if r.msgMap == nil {
		r.msgMap = map[int]MsgParser{}
	}
	if r.typeMap == nil {
		r.typeMap = map[reflect.Type]*MessageID{}
	}
	mId := &MessageID{Cmd: cmd, Act: act}
	p := MsgParser{}
	if c2s != nil {
		c2sType := reflect.TypeOf(c2s).Elem()
		p.c2sFunc = func() interface{} {
			return reflect.New(c2sType).Interface()
		}
		if v, ok := r.typeMap[c2sType]; ok {
			panic(Sprintf("Parser register repeat c2sType: %v, cmd: %v, act: %v, old: %v", c2sType, uint8(cmd), uint8(act), v))
		}
		r.typeMap[c2sType] = mId
	}
	if s2c != nil {
		s2cType := reflect.TypeOf(s2c).Elem()
		p.s2cFunc = func() interface{} {
			return reflect.New(s2cType).Interface()
		}
		if v, ok := r.typeMap[s2cType]; ok {
			panic(Sprintf("Parser register repeat s2cType: %v, cmd: %v, act: %v, old: %v", s2cType, uint8(cmd), uint8(act), v))
		}
		r.typeMap[s2cType] = mId
		//LogInfo("==============%v, %v, %v", s2cType, uint8(cmd), uint8(act))
	}
	if (cmd < 0 || cmd > 255) || (act < 0 || act > 255) {
		panic(Sprintf("Parser register cmd act error: %v, %v", cmd, act))
	}
	msgId := mId.ID()
	if parser, ok := r.msgMap[msgId]; ok {
		panic(Sprintf("Parser register repeat cmd: %v, act: %v, parser: %v, %v", uint8(cmd), uint8(act), reflect.TypeOf(parser.C2S()), reflect.TypeOf(p.C2S())))
	}
	r.msgMap[msgId] = p
}

func (r *Parser) RegisterMsgFunc(c2sFunc ParseFunc, s2cFunc ParseFunc) {
	if r.cmdRoot == nil {
		r.cmdRoot = &cmdParseNode{}
	}
	registerCmdParser(r.cmdRoot, c2sFunc, s2cFunc)
}

func (r *Parser) RegisterMsg(c2s interface{}, s2c interface{}) {
	var c2sFunc ParseFunc = nil
	var s2cFunc ParseFunc = nil
	if c2s != nil {
		c2sType := reflect.TypeOf(c2s).Elem()
		c2sFunc = func() interface{} {
			return reflect.New(c2sType).Interface()
		}
	}
	if s2c != nil {
		s2cType := reflect.TypeOf(s2c).Elem()
		s2cFunc = func() interface{} {
			return reflect.New(s2cType).Interface()
		}
	}

	if r.cmdRoot == nil {
		r.cmdRoot = &cmdParseNode{}
	}
	registerCmdParser(r.cmdRoot, c2sFunc, s2cFunc)
}

func JsonUnPack(data []byte, msg interface{}) error {
	if data == nil || msg == nil {
		return ErrJsonUnPack
	}

	err := json.Unmarshal(data, msg)
	if err != nil {
		LogInfo("unmarshal json err:%v data:%s", err, data)
		return ErrJsonUnPack
	}
	return nil
}

func JsonPack(msg interface{}) ([]byte, error) {
	if msg == nil {
		return nil, ErrJsonPack
	}

	data, err := json.Marshal(msg)
	if err != nil {
		LogInfo("marshal json:%v err:%v", msg, err)
		return nil, ErrJsonPack
	}

	return data, nil
}

func MsgPackUnPack(data []byte, msg interface{}) error {
	err := msgpack.Unmarshal(data, msg)
	return err
}

func MsgPackPack(msg interface{}) ([]byte, error) {
	data, err := msgpack.Marshal(msg)
	return data, err
}
