package antnet

import (
	"github.com/gogo/protobuf/proto"
)

type pBParser struct {
	*Parser
}

func (r *pBParser) ParseC2S(msg *Message) (IMsgParser, error) {
	if msg == nil {
		LogError("ParseC2S message is nil")
		return nil, ErrPBUnPack
	}
	// 某些消息不用解析
	if msg.Head.Flags&FlagNoParse > 0 {
		return nil, nil
	}
	p, ok := r.msgMap[msg.Head.CmdAct()]
	if ok {
		if p.C2S() != nil {
			err := PBUnPack(msg.Data, p.C2S())
			if err != nil {
				dLen := len(msg.Data)
				var text []byte
				if dLen > 100 {
					text = msg.Data[:100]
				} else {
					text = msg.Data
				}
				LogError("ParseC2S PBUnPack error cmd: %v, act: %v, err: %v, len: %d/%d, data: %v", msg.Head.Cmd, msg.Head.Act, err, len(text), dLen, text)
				return nil, err
			}
			p.parser = r
			return &p, nil
		} else {
			return &p, nil
		}
	} else {
		LogError("ParseC2S not register message cmd: %v, act: %v", msg.Head.Cmd, msg.Head.Act)
		r.ErrType = ParseErrTypeContinue
	}

	return nil, ErrPBUnPack
}

func (r *pBParser) PackMsg(v interface{}) []byte {
	data, _ := PBPack(v)
	return data
}

func PBUnPack(data []byte, msg interface{}) error {
	//if data == nil || msg == nil {
	if msg == nil {
		return ErrPBUnPack
	}

	err := proto.Unmarshal(data, msg.(proto.Message))
	if err != nil {
		dLen := len(data)
		var text []byte
		if dLen > 100 {
			text = data[:100]
		} else {
			text = data
		}
		LogWarn("PBUnPack Unmarshal error: %s, msg type: %v, len: %d/%d, data: %v", err.Error(), basal.Type(msg), len(text), dLen, text)
		return ErrPBUnPack
	}
	return nil
}

func PBPack(msg interface{}) ([]byte, error) {
	if msg == nil {
		return nil, ErrPBPack
	}

	data, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		LogError("PBPack Marshal error:%s, msg: %v", err, msg)
	}

	return data, nil
}
