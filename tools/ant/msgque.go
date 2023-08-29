package antnet

import (
	"core/tools/core"
	"encoding/binary"
	"github.com/xtaci/kcp-go"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

var DefMsgQueTimeout = 180

type MsgType int

const (
	MsgTypeMsg MsgType = iota //消息基于确定的消息头
	MsgTypeCmd                //消息没有消息头，以\n分割
	MsgTypeRpc
)

type NetType int

const (
	NetTypeTcp NetType = iota //TCP类型
	NetTypeUdp                //UDP类型
	NetTypeWs                 //websocket
)

type ConnType int

const (
	ConnTypeListen ConnType = iota //监听
	ConnTypeConn                   //连接产生的
	ConnTypeAccept                 //Accept产生的
)

type IMsgQue interface {
	Id() uint32
	GetMsgType() MsgType
	GetConnType() ConnType
	GetNetType() NetType
	Wait()

	LocalAddr() string
	RemoteAddr() string
	RemoteIP() string

	Stop()
	IsStop() bool
	Available() bool

	Send(m *Message) (re bool)
	SendString(str string) (re bool)
	SendStringLn(str string) (re bool)
	SendByteStr(str []byte) (re bool)
	SendByteStrLn(str []byte) (re bool)
	SendCallback(m *Message, c chan *Message) (re bool)
	SetSendFast()
	SetTimeout(t int)
	GetTimeout() int
	Reconnect(t int) //重连间隔  最小1s，此函数仅能连接关闭是调用

	GetHandler() IMsgHandler

	SetUser(user interface{})
	GetUser() interface{}

	Store(key, value interface{})
	Load(key interface{}) (value interface{}, ok bool)
	Delete(key interface{})

	tryCallback(msg *Message) (re bool)
	SendLen() int
}

type msgQue struct {
	sync.Map

	id uint32 //唯一标示

	cwrite  chan *Message //写入通道
	stop    int32         //停止标记
	msgTyp  MsgType       //消息类型
	connTyp ConnType      //通道类型

	handler       IMsgHandler //处理者
	parser        IParser
	parserFactory *Parser
	timeout       int //传输超时
	lastTick      int64

	init         bool
	available    bool
	sendFast     bool
	callback     map[uint16]chan *Message
	user         interface{}
	callbackLock sync.Mutex
	//gmsgId       uint16

	encrypt bool   //通信是否加密
	iseed   uint32 //input种子
	oseed   uint32 //output种子

	waitRecvSeed chan struct{} //是否接收到种子
}

func (r *msgQue) SendLen() int {
	defer core.Exception(nil)
	return len(r.cwrite)
}

func (r *msgQue) SetSendFast() {
	r.sendFast = true
}

func (r *msgQue) SetUser(user interface{}) {
	r.user = user
}

func (r *msgQue) SetEncrypt(e bool) {
	r.encrypt = e
	if e && r.connTyp == ConnTypeAccept {
		r.oseed = uint32(Timestamp)
		r.iseed = uint32(Timestamp) + uint32(RandNumber(99999))

		data := make([]byte, 8)
		binary.BigEndian.PutUint32(data, r.iseed)
		binary.BigEndian.PutUint32(data[4:], r.oseed)
		msg := NewMsg(0, 0, data)
		//LogInfo("init session msg:%v iseed:%d oseed:%d", msg, r.iseed, r.oseed)
		r.cwrite <- msg
	}
}

func (r *msgQue) GetEncrypt() bool {
	return r.encrypt
}

//func (r *msgQue) getGMsg(add bool) *gMsg {
//	if add {
//		r.gmsgId++
//	}
//	gm := gmsgArray[r.gmsgId]
//	return gm
//}

func (r *msgQue) Available() bool {
	return r.available
}

func (r *msgQue) GetUser() interface{} {
	return r.user
}

func (r *msgQue) GetHandler() IMsgHandler {
	return r.handler
}

func (r *msgQue) GetMsgType() MsgType {
	return r.msgTyp
}

func (r *msgQue) GetConnType() ConnType {
	return r.connTyp
}

func (r *msgQue) Id() uint32 {
	return r.id
}

func (r *msgQue) SetTimeout(t int) {
	if t >= 0 {
		r.timeout = t
	}
}

func (r *msgQue) isTimeout(tick *time.Timer) bool {
	left := int(Timestamp - r.lastTick)
	if left < r.timeout || r.timeout == 0 {
		if r.timeout == 0 {
			tick.Reset(time.Second * time.Duration(DefMsgQueTimeout))
		} else {
			tick.Reset(time.Second * time.Duration(r.timeout-left))
		}
		return false
	}
	LogInfo("msgque close because timeout id:%v wait:%v timeout:%v user:%v", r.id, left, r.timeout, r.user)
	return true
}

func (r *msgQue) GetTimeout() int {
	return r.timeout
}

func (r *msgQue) Reconnect(t int) {

}

func (r *msgQue) Send(m *Message) (re bool) {
	if m == nil || !r.available {
		return
	}
	//if len(r.cwrite) > cwrite_chan_len-1 {
	//	return
	//}
	defer func() {
		if err := recover(); err != nil {
			re = false
		}
	}()
	//LogInfo("11111111111111111 :%v, %v", m.Head.Len, m.Data)
	//if m.Head != nil && m.Head.Len >= 1500 && (m.Head.Flags&FlagCompress) == 0 {
	//	//LogInfo("Start GZipCompress head len:%v [%v %v]", m.Head.Len, m.Head.Cmd, m.Head.Act)
	//	m.Head.Flags |= FlagCompress
	//	m.Data = GZipCompress(m.Data)
	//	m.Head.Len = uint32(len(m.Data))
	//	//LogInfo("End GZipCompress 2 head len:%v", m.Head.Len)
	//}

	//if r.encrypt && m.Head != nil {
	//	m = m.Copy()
	//	//LogInfo("Start Encrypt head:%v data:%v", m.Head, m.Data)
	//	m.Head.Flags |= FlagEncrypt
	//	r.oseed = r.oseed*cryptA + cryptB
	//	m.Head.Bcc = CountBCC(m.Data, 0, m.Head.Len)
	//	m.Data = DefaultNetEncrypt(r.oseed, m.Data, 0, m.Head.Len)
	//	//LogInfo("End Encrypt seed:%d bcc:%v head:%v data:%v", r.oseed, m.Head.Bcc, m.Head, m.Data)
	//}

	//LogInfo("msgQue send 11  len:%v", len(r.cwrite))
	r.cwrite <- m
	//LogInfo("msgQue send 22  len:%v", len(r.cwrite))
	return true
}

func (r *msgQue) SendCallback(m *Message, c chan *Message) (re bool) {
	if c == nil || cap(c) < 1 {
		LogError("try send callback but chan is null or no buffer")
		return
	}
	if r.Send(m) {
		r.setCallback(m.Tag(), c)
	} else {
		c <- nil
		return
	}
	return true
}

func (r *msgQue) SendString(str string) (re bool) {
	return r.Send(&Message{Data: []byte(str)})
}

func (r *msgQue) SendStringLn(str string) (re bool) {
	return r.SendString(str + "\n")
}

func (r *msgQue) SendByteStr(str []byte) (re bool) {
	return r.SendString(string(str))
}

func (r *msgQue) SendByteStrLn(str []byte) (re bool) {
	return r.SendString(string(str) + "\n")
}

func (r *msgQue) tryCallback(msg *Message) (re bool) {
	if r.callback == nil {
		return false
	}
	defer func() {
		if err := recover(); err != nil {

		}
		r.callbackLock.Unlock()
	}()
	r.callbackLock.Lock()
	if r.callback != nil {
		tag := msg.Tag()
		if c, ok := r.callback[tag]; ok {
			delete(r.callback, tag)
			c <- msg
			re = true
		}
	}
	return
}

func (r *msgQue) setCallback(tag uint16, c chan *Message) {
	defer func() {
		if err := recover(); err != nil {

		}
		r.callback[tag] = c
		r.callbackLock.Unlock()
	}()

	r.callbackLock.Lock()
	if r.callback == nil {
		r.callback = make(map[uint16]chan *Message)
	}
	oc, ok := r.callback[tag]
	if ok { //可能已经关闭
		oc <- nil
	}
}

func (r *msgQue) baseStop() {
	if r.cwrite != nil {
		close(r.cwrite)
	}

	for k, v := range r.callback {
		v <- nil
		delete(r.callback, k)
	}
	msgqueMapSync.Lock()
	delete(msgqueMap, r.id)
	msgqueMapSync.Unlock()
	LogInfo("msgque close id:%d", r.id)
}

func (r *msgQue) processMsg(msgque IMsgQue, msg *Message) bool {
	//作为客户端时--robot
	if r.connTyp == ConnTypeConn && msg.Head.Cmd == 0 && msg.Head.Act == 0 {
		if len(msg.Data) != 9 {
			LogWarn("init seed msg err: %d", len(msg.Data))
			return false
		}
		r.encrypt = true
		msg.RemoveDataPlaceholder()
		r.oseed = binary.BigEndian.Uint32(msg.Data[:4])
		r.iseed = binary.BigEndian.Uint32(msg.Data[4:])
		if r.waitRecvSeed != nil {
			core.Try(func() {
				close(r.waitRecvSeed)
				//LogInfo("收到加密种子: %v", r.id)
			}, nil)
		}
		return true
	}

	if msg.Head != nil && msg.Head.Flags&FlagEncrypt > 0 {
		//if msg.Head.Flags&FlagEncrypt <= 0 {
		//	LogError("Decrypt head encrypt true,flags error:%v", msg.Head)
		//	return false
		//}
		//LogInfo("Start Decrypt head:%v data:%v", msg.Head, msg.Data)
		//LogInfo("Start Decrypt seed:%d ", r.iseed)
		r.iseed = r.iseed*cryptA + cryptB
		msg.Data = DefaultNetDecrypt(r.iseed, msg.Data, 0, msg.Head.Len)
		bcc := CountBCC(msg.Data, 0, msg.Head.Len)
		//LogInfo("End Decrypt seed:%d bcc:%v Head:%v data:%v", r.iseed, bcc, msg.Head, msg.Data)
		if msg.Head.Bcc != bcc {
			LogWarn("client bcc err conn:%d, bcc: %v, head bcc: %v", r.id, bcc, msg.Head.Bcc)
			return false
		}
	}

	if msg.Head != nil && msg.Head.Flags&FlagCompress > 0 && msg.Data != nil {
		data, err := GZipUnCompress(msg.Data)
		if err != nil {
			LogError("msgque uncompress failed msgque:%v cmd:%v act:%v len:%v err:%v", msgque.Id(), msg.Head.Cmd, msg.Head.Act, msg.Head.Len, err)
			return false
		}
		msg.Data = data
		msg.Head.Len = uint32(len(msg.Data))
	}

	if r.encrypt {
		//LogInfo("000000000000=========cmd: %d, act: %d, len: %d, data: %v", msg.Cmd(), msg.Act(), msg.Len(), msg.Data)
		if !msg.RemoveDataPlaceholder() {
			LogError("msgque RemoveDataPlaceholder failed msgque:%v cmd:%v act:%v len:%v", msgque.Id(), msg.Head.Cmd, msg.Head.Act, msg.Head.Len)
			return false
		}
		//LogInfo("111111111111=========cmd: %d, act: %d, len: %d, data: %v", msg.Cmd(), msg.Act(), msg.Len(), msg.Data)
	}

	if r.parser != nil {
		mp, err := r.parser.ParseC2S(msg)
		if err == nil {
			msg.IMsgParser = mp
		} else {
			if r.parser.GetErrType() == ParseErrTypeSendRemind {
				//机器人作为客户端时，不能全部注册消息，但不能断开连接
				if r.connTyp == ConnTypeConn {
					return true
				}
				LogError("parse msg error: cmd: %v, act: %v, addr: %v, err: %v", msg.Cmd(), msg.Act(), msgque.RemoteAddr(), err)
				return false
			} else if r.parser.GetErrType() == ParseErrTypeClose {
				LogError("parse msg error ParseErrTypeClose: cmd: %v, act: %v, addr: %v, err: %v", msg.Cmd(), msg.Act(), msgque.RemoteAddr(), err)
				return false
			} else if r.parser.GetErrType() == ParseErrTypeContinue {
				LogError("parse msg error ParseErrTypeContinue: cmd: %v, act: %v, addr: %v, err: %v", msg.Cmd(), msg.Act(), msgque.RemoteAddr(), err)
				return true
			}
		}
	}
	f := r.handler.GetHandlerFunc(msgque, msg)
	if f == nil {
		f = r.handler.OnProcessMsg
	}

	var status bool
	Try(func() {
		status = f(msgque, msg)
	}, func(stack string, e error) {
		status = true
	})
	return status
}

type HandlerFunc func(msgque IMsgQue, msg *Message) bool

type IMsgHandler interface {
	OnNewMsgQue(msgque IMsgQue) bool                         //新的消息队列
	OnDelMsgQue(msgque IMsgQue)                              //消息队列关闭
	OnProcessMsg(msgque IMsgQue, msg *Message) bool          //默认的消息处理函数
	OnConnectComplete(msgque IMsgQue, ok bool) bool          //连接成功
	GetHandlerFunc(msgque IMsgQue, msg *Message) HandlerFunc //根据消息获得处理函数
}

type IMsgRegister interface {
	Register(cmd, act uint8, fun HandlerFunc)
	RegisterMsg(v interface{}, fun HandlerFunc)
}

type DefMsgHandler struct {
	msgMap  map[int]HandlerFunc
	typeMap map[reflect.Type]HandlerFunc
}

func (r *DefMsgHandler) OnNewMsgQue(msgque IMsgQue) bool                { return true }
func (r *DefMsgHandler) OnDelMsgQue(msgque IMsgQue)                     {}
func (r *DefMsgHandler) OnProcessMsg(msgque IMsgQue, msg *Message) bool { return true }
func (r *DefMsgHandler) OnConnectComplete(msgque IMsgQue, ok bool) bool { return true }
func (r *DefMsgHandler) GetHandlerFunc(msgque IMsgQue, msg *Message) HandlerFunc {
	if msgque.tryCallback(msg) {
		return r.OnProcessMsg
	}

	if msg.CmdAct() == 0 {
		if r.typeMap != nil {
			if f, ok := r.typeMap[reflect.TypeOf(msg.C2S())]; ok {
				return f
			}
		}
	} else if r.msgMap != nil {
		if f, ok := r.msgMap[msg.CmdAct()]; ok {
			return f
		}
	}

	return nil
}

func (r *DefMsgHandler) RegisterMsg(v interface{}, fun HandlerFunc) {
	msgType := reflect.TypeOf(v)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		LogFatal("message pointer required")
		return
	}
	if r.typeMap == nil {
		r.typeMap = map[reflect.Type]HandlerFunc{}
	}
	r.typeMap[msgType] = fun
}

func (r *DefMsgHandler) Register(cmd uint8, act uint8, fun HandlerFunc) {
	if r.msgMap == nil {
		r.msgMap = map[int]HandlerFunc{}
	}
	if (cmd < 0 || cmd > 255) || (act < 0 || act > 255) {
		panic(Sprintf("DefMsgHandler register cmd act error: %v, %v", cmd, act))
	}
	r.msgMap[CmdAct(cmd, act)] = fun
}

type EchoMsgHandler struct {
	DefMsgHandler
}

func (r *EchoMsgHandler) OnProcessMsg(msgque IMsgQue, msg *Message) bool {
	msgque.Send(msg)
	return true
}

func StartServer(addr string, typ MsgType, handler IMsgHandler, parser *Parser, encrypt bool) ([]IMsgQue, error) {
	msgques := make([]IMsgQue, 0, 1)
	addrs := strings.Split(addr, "://")
	if addrs[0] == "tcp" || addrs[0] == "all" {
		listen, err := net.Listen("tcp", addrs[1])
		if err == nil {
			msgque := newTcpListen(listen, typ, handler, parser, addr)
			msgque.SetEncrypt(encrypt)
			Go(func() {
				LogDebug("process listen for msgque:%d", msgque.id)
				msgque.listen()
				LogDebug("process listen end for msgque:%d", msgque.id)
			})
			msgques = append(msgques, msgque)
		} else {
			LogError("listen on %s failed, errstr:%s", addr, err)
			return nil, err
		}
	}
	//if addrs[0] == "udp" || addrs[0] == "all" {
	//	naddr, err := net.ResolveUDPAddr("udp", addrs[1])
	//	if err != nil {
	//		LogError("listen on %s failed, errstr:%s", addr, err)
	//		return err
	//	}
	//	conn, err := net.ListenUDP("udp", naddr)
	//	if err == nil {
	//		msgque := newUdpListen(conn, typ, handler, parser, addr)
	//		Go(func() {
	//			LogDebug("process listen for msgque:%d", msgque.id)
	//			msgque.listen()
	//			LogDebug("process listen end for msgque:%d", msgque.id)
	//		})
	//	} else {
	//		LogError("listen on %s failed, errstr:%s", addr, err)
	//		return err
	//	}
	//}
	if addrs[0] == "udp" || addrs[0] == "all" {
		lis, err := kcp.ListenWithOptions(addrs[1], nil, 0, 0)
		if err == nil {
			LogInfo("listening on:%v", lis.Addr())
			msgque := newKcpListen(lis, typ, handler, parser, addr)
			msgque.SetEncrypt(encrypt)
			Go(func() {
				LogDebug("process listen for msgque:%d", msgque.id)
				msgque.listen()
				LogDebug("process listen end for msgque:%d", msgque.id)
			})
			msgques = append(msgques, msgque)
		}
	}
	return msgques, nil
}

func StartConnect(netype string, addr string, typ MsgType, handler IMsgHandler, parser *Parser, user interface{}, encrypt bool) IMsgQue {
	if IsStop() {
		return nil
	}
	if netype == "tcp" {
		msgque := newTcpConn(netype, addr, nil, typ, handler, parser, user)
		msgque.encrypt = encrypt
		if handler.OnNewMsgQue(msgque) {
			msgque.init = true
			if msgque.Connect() {
				return msgque
			}
			LogError("connect to:%s:%s failed", netype, addr)
		} else {
			msgque.Stop()
		}
	} else if netype == "udp" || netype == "kcp" {
		msgque := newKcpConn(netype, addr, nil, typ, handler, parser, user)
		if handler.OnNewMsgQue(msgque) {
			msgque.init = true
			if msgque.Connect() {
				return msgque
			}
			LogError("connect to:%s:%s failed", netype, addr)
		} else {
			msgque.Stop()
		}
	}

	return nil
}
