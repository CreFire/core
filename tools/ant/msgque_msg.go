package antnet

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"unsafe"
)

const (
	MsgHeadSize = 8
)

const (
	FlagEncrypt  = 1 << 0 //数据是经过加密的
	FlagCompress = 1 << 1 //数据是经过压缩的
	FlagContinue = 1 << 2 //消息还有后续
	FlagNeedAck  = 1 << 3 //消息需要确认
	FlagAck      = 1 << 4 //确认消息
	FlagReSend   = 1 << 5 //重发消息
	FlagClient   = 1 << 6 //消息来自客服端，用于判断index来之服务器还是其他玩家
	FlagNoParse  = 1 << 7 //消息不解析
)

const MaxMsgDataSize uint32 = 1024 * 1024 //收到的数据最大长度

type MessageHead struct {
	Len   uint32 //数据长度
	Cmd   uint8  //命令
	Act   uint8  //动作
	Flags uint8  //标记
	Bcc   uint8  //加密校验
	data  []byte
}

func (r *MessageHead) CMD() uint8 {
	return uint8(r.Cmd)
}

func (r *MessageHead) ACT() uint8 {
	return uint8(r.Act)
}

func (r *MessageHead) Copy() *MessageHead {
	if r == nil {
		return nil
	}
	head := &MessageHead{
		Len:   r.Len,
		Cmd:   r.Cmd,
		Act:   r.Act,
		Flags: r.Flags,
		Bcc:   r.Bcc,
	}
	ld := len(r.data)
	if ld > 0 {
		head.data = make([]byte, ld)
		copy(head.data, r.data)
	}
	return head
}

func (r *MessageHead) Bytes() []byte {
	if r.data != nil {
		return r.data
	}
	r.data = make([]byte, MsgHeadSize)
	phead := (*MessageHead)(unsafe.Pointer(&r.data[0]))
	phead.Len = r.Len
	phead.Cmd = r.Cmd
	phead.Act = r.Act
	phead.Flags = r.Flags
	phead.Bcc = r.Bcc
	return r.data
}

func (r *MessageHead) FastBytes(data []byte) []byte {
	phead := (*MessageHead)(unsafe.Pointer(&data[0]))
	phead.Len = r.Len
	phead.Cmd = r.Cmd
	phead.Act = r.Act
	phead.Flags = r.Flags
	phead.Bcc = r.Bcc
	return data
}

func (r *MessageHead) BytesWithData(wdata []byte) []byte {
	if r.data != nil {
		return r.data
	}
	r.Len = uint32(len(wdata))
	r.data = make([]byte, MsgHeadSize+r.Len)
	phead := (*MessageHead)(unsafe.Pointer(&r.data[0]))
	phead.Len = r.Len
	phead.Cmd = r.Cmd
	phead.Act = r.Act
	phead.Flags = r.Flags
	phead.Bcc = r.Bcc
	if wdata != nil {
		copy(r.data[MsgHeadSize:], wdata)
	}
	return r.data
}

func (r *MessageHead) FromBytes(data []byte) error {
	if len(data) < MsgHeadSize {
		return ErrMsgLenTooShort
	}
	phead := (*MessageHead)(unsafe.Pointer(&data[0]))
	r.Len = phead.Len
	r.Cmd = phead.Cmd
	r.Act = phead.Act
	r.Flags = phead.Flags
	r.Bcc = phead.Bcc
	if r.Len > MaxMsgDataSize {
		return ErrMsgLenTooLong
	}
	return nil
}

func (r *MessageHead) CmdAct() int {
	return CmdAct(r.Cmd, r.Act)
}

func (r *MessageHead) Tag() uint16 {
	return Tag(r.Cmd, r.Act)
}

func (r *MessageHead) String() string {
	return fmt.Sprintf("Len:%v Cmd:%v Act:%v Flags:%v", r.Len, r.Cmd, r.Act, r.Flags)
}

func NewMessageHead(data []byte) *MessageHead {
	head := &MessageHead{}
	if err := head.FromBytes(data); err != nil {
		LogError("NewMessageHead error: %v", err)
		return nil
	}
	return head
}

func MessageHeadFromByte(data []byte) *MessageHead {
	if len(data) < MsgHeadSize {
		return nil
	}
	phead := new(*MessageHead)
	*phead = (*MessageHead)(unsafe.Pointer(&data[0]))
	if (*phead).Len > MaxMsgDataSize {
		return nil
	}
	return *phead
}

type Message struct {
	Head         *MessageHead //消息头，可能为nil
	Data         []byte       //消息数据
	IMsgParser                //解析器
	User         interface{}  //用户自定义数据
	serializeC2S bool
}

// 增加占位符
func (r *Message) AddDataPlaceholder() bool {
	if r.Head != nil {
		dLen := len(r.Data) + 1
		newData := make([]byte, dLen)
		copy(newData[1:], r.Data)
		r.Data = newData
		r.Head.Len = uint32(dLen)
		return true
	} else {
		LogError("AddDataPlaceholder failed")
		return false
	}
}

// 移除占位符
func (r *Message) RemoveDataPlaceholder() bool {
	if r.Head != nil {
		if len(r.Data) > 0 {
			r.Data = r.Data[1:]
			r.Head.Len = uint32(len(r.Data))
			return true
		} else {
			LogError("RemoveDataPlaceholder failed 1")
			return false
		}
	} else {
		LogError("RemoveDataPlaceholder failed 2")
		return false
	}
}

func (r *Message) Copy() *Message {
	msg := &Message{
		Head:       r.Head.Copy(),
		User:       r.User,
		IMsgParser: r.IMsgParser,
	}
	ld := len(r.Data)
	if ld > 0 {
		msg.Data = make([]byte, ld)
		copy(msg.Data, r.Data)
	}
	return msg
}

func (r *Message) CmdAct() int {
	if r.Head != nil {
		return CmdAct(r.Head.Cmd, r.Head.Act)
	}
	return 0
}

func (r *Message) Len() uint32 {
	if r.Head != nil {
		return r.Head.Len
	}
	return 0
}

func (r *Message) Cmd() uint8 {
	if r.Head != nil {
		return r.Head.Cmd
	}
	return 0
}

func (r *Message) Act() uint8 {
	if r.Head != nil {
		return r.Head.Act
	}
	return 0
}

func (r *Message) Flags() uint8 {
	if r.Head != nil {
		return r.Head.Flags
	}
	return 0
}

func (r *Message) Bcc() uint8 {
	if r.Head != nil {
		return r.Head.Bcc
	}
	return 0
}

func (r *Message) Tag() uint16 {
	if r.Head != nil {
		return Tag(r.Head.Cmd, r.Head.Act)
	}
	return 0
}

func (r *Message) Bytes() []byte {
	if r.Head != nil {
		if r.Data != nil {
			return r.Head.BytesWithData(r.Data)
		}
		return r.Head.Bytes()
	}
	return r.Data
}

func (r *Message) CopyTag(old *Message) *Message {
	if r.Head != nil && old.Head != nil {
		r.Head.Cmd = old.Head.Cmd
		r.Head.Act = old.Head.Act
	}
	return r
}

func (r *Message) SerializeC2S() (state bool) {
	Try(func() {
		if r.serializeC2S {
			c2s := r.C2S()
			if c2s == nil {
				panic(Sprintf("SerializeC2S c2s nil: %v", c2s))
			}
			var err error
			r.Data, err = PBPack(c2s)
			if err != nil {
				panic(Sprintf("SerializeC2S PBPack error: %v, %v", err, c2s))
			}
			r.Head.Len = uint32(len(r.Data))
			state = true
		} else {
			state = false
		}
	}, nil)
	return
}

func NewStrMsg(str string) *Message {
	return &Message{
		Data: []byte(str),
	}
}

func NewPbMsg(cmd uint8, act uint8, msg proto.Message, parser IMsgParser) *Message {
	//data, err := PBPack(msg)
	//if err != nil {
	//	panic(Sprintf("NewPbMsg PBPack error: %v, msg: %v", err, msg))
	//}
	parser.setC2S(msg)
	pbMsg := &Message{
		Head: &MessageHead{
			//Len: uint32(len(data)),
			Cmd: uint8(cmd),
			Act: uint8(act),
		},
		//Data:       data,
		IMsgParser:   parser,
		serializeC2S: true,
	}
	return pbMsg
}

func NewDataMsg(data []byte) *Message {
	return &Message{
		Head: &MessageHead{
			Len: uint32(len(data)),
		},
		Data: data,
	}
}

func NewMsg(cmd uint8, act uint8, data []byte) *Message {

	return &Message{
		Head: &MessageHead{
			Len: uint32(len(data)),
			Cmd: uint8(cmd),
			Act: uint8(act),
		},
		Data: data,
	}

}

func NewTagMsg(cmd uint8, act uint8) *Message {
	return &Message{
		Head: &MessageHead{
			Cmd: uint8(cmd),
			Act: uint8(act),
		},
	}
}
