package antnet

import (
	"github.com/xtaci/kcp-go"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var cwrite_chan_len = 64
var cread_chan_len = 64

type kcpMsgQue struct {
	msgQue
	lis         *kcp.Listener
	session     *kcp.UDPSession
	cread       chan []byte //写入通道
	addr        net.Addr
	network     string
	address     string
	connecting  int32
	chReadEvent chan struct{} // notify Read() can be called without blocking
	sync.Mutex
}

func (r *kcpMsgQue) GetNetType() NetType {
	return NetTypeUdp
}

func (r *kcpMsgQue) Wait() {

}

func (r *kcpMsgQue) Stop() {
	LogInfo("kcpMsgQue stop")
	if atomic.CompareAndSwapInt32(&r.stop, 0, 1) {
		Go(func() {
			if r.init {
				r.handler.OnDelMsgQue(r)
			}
			r.available = false
			if r.cread != nil {
				close(r.cread)
			}
			if r.session != nil {
				r.session.Close()
			}
			r.baseStop()
		})
	}
}

func (r *kcpMsgQue) IsStop() bool {
	if r.stop == 0 {
		if IsStop() {
			LogError("kcpMsgQue IsStop stop")
			r.Stop()
		}
	}
	return r.stop == 1
}

func (r *kcpMsgQue) LocalAddr() string {
	if r.addr != nil {
		return r.addr.String()
	}
	return ""
}

func (r *kcpMsgQue) RemoteAddr() string {
	if r.session != nil {
		return r.session.RemoteAddr().String()
	}
	return ""
}

func (r *kcpMsgQue) RemoteIP() string {
	addr := strings.Split(r.RemoteAddr(), ":")
	return addr[0]
}

func (r *kcpMsgQue) read() {
	defer func() {
		if err := recover(); err != nil {
			LogError("msgque read panic id:%v err:%v", r.id, err.(error))
			LogStack()
		}
		r.Stop()
	}()
	data := make([]byte, 1<<16)
	for !r.IsStop() {
		n, _ := r.session.Read(data)
		if data == nil {
			LogError("kcpMsgQue stop Read data nil")
			break
		}
		pdata := make([]byte, n)
		copy(pdata, data)
		//fmt.Printf("<<<<<<<<<<<<<<<<<kcp 33 recev:%v \n", pdata)
		var msg *Message
		if r.msgTyp == MsgTypeCmd {
			msg = &Message{Data: pdata}
		} else {
			head := MessageHeadFromByte(pdata)
			if head == nil {
				LogError("kcpMsgQue stop head nil")
				break
			}
			if head.Len > 0 {
				msg = &Message{Head: head, Data: pdata[MsgHeadSize:]}
			} else {
				msg = &Message{Head: head}
			}
		}
		r.lastTick = Timestamp
		if !r.init {
			//fmt.Printf("<<<<<<<<<<<<<<<<<kcp 44 recev:%v \n", pdata)
			if !r.handler.OnNewMsgQue(r) {
				LogError("kcpMsgQue stop OnNewMsgQue error")
				break
			}
			r.init = true
		}
		//fmt.Printf("<<<<<<<<<<<<<<<<<kcp 55 recev msghead:%v msgdata:%v \n", msg.Head, msg.Data)
		if !r.processMsg(r, msg) {
			LogError("kcpMsgQue stop processMsg error")
			break
		}
	}
	r.Stop()
}

func (r *kcpMsgQue) write() {
	//fmt.Printf("<<<<<<<<<<<<<<<<<kcp 00 write \n")
	defer func() {
		if err := recover(); err != nil {
			LogError("msgque write panic id:%v err:%v", r.id, err.(error))
			LogStack()
		}
		r.Stop()
	}()
	//gm := r.getGMsg(false)
	tick := time.NewTimer(time.Second * time.Duration(r.timeout))
	for !r.IsStop() {
		//fmt.Printf("<<<<<<<<<<<<<<<<<kcp 11 write \n")
		var m *Message = nil
		select {
		case <-stopChanForGo:
		case m = <-r.cwrite:
		//case <-gm.c:
		//	if gm.fun == nil || gm.fun(r) {
		//		m = gm.msg
		//	}
		//	gm = r.getGMsg(true)
		case <-tick.C:
			if r.isTimeout(tick) {
				//fmt.Printf("<<<<<<<<<<<<<<<<<kcp write stop timeout:%v \n", r.timeout)
				r.Stop()
			}
		}
		//fmt.Printf("<<<<<<<<<<<<<<<<<kcp 11 write:%v \n", m)
		if m == nil {
			//fmt.Printf("<<<<<<<<<<<<<<<<<kcp m nil \n")
			continue
		}

		if r.msgTyp == MsgTypeCmd {
			if m.Data != nil {
				//r.conn.WriteToUDP(m.Data, r.addr)
				r.session.Write(m.Data)
			}
		} else {
			if m.Head != nil || m.Data != nil {
				//r.conn.WriteToUDP(m.Bytes(), r.addr)
				//fmt.Printf("<<<<<<<<<<<<<<<<<kcp 22 write \n")
				r.session.Write(m.Bytes())
				//fmt.Printf("<<<<<<<<<<<<<<<<<kcp 33 write \n")
			}
		}

		r.lastTick = Timestamp
	}
	tick.Stop()
	r.Stop()
}

func (r *kcpMsgQue) listen() {
	Go2(func(cstop chan struct{}) {
		for {
			select {
			case <-cstop:
				r.lis.Close()
				return
			}
		}
	})

	for !r.IsStop() {
		session, err := r.lis.AcceptKCP()
		if err != nil {
			if stop == 0 && r.stop == 0 {
				LogError("accept failed msgque:%v err:%v", r.id, err)
			}
			break
		} else {
			session.SetStreamMode(false)
			session.SetWriteDelay(false)
			session.SetNoDelay(1, 10, 2, 1)
			//session.SetMtu(500)
			session.SetWindowSize(512, 512)
			session.SetACKNoDelay(true)
			Go(func() {
				msgque := newKcpAccept(session, r.msgTyp, r.handler, r.parserFactory)
				msgque.SetEncrypt(r.GetEncrypt())
				Go(func() {
					LogInfo("process kcp read for msgque:%d", msgque.id)
					msgque.read()
					LogInfo("process kcp read end for msgque:%d", msgque.id)
				})
				Go(func() {
					LogInfo("process kcp write for msgque:%d", msgque.id)
					msgque.write()
					LogInfo("process kcp write end for msgque:%d", msgque.id)
				})
			})
		}
	}
}

func (r *kcpMsgQue) Connect() bool {
	LogInfo("connect to addr:%s msgque:%d start", r.address, r.id)
	c, err := kcp.DialWithOptions(r.address, nil, 0, 0)
	if err != nil {
		LogInfo("connect to addr:%s msgque:%d err:%v", r.address, r.id, err)
		r.handler.OnConnectComplete(r, false)
		atomic.CompareAndSwapInt32(&r.connecting, 1, 0)
		r.Stop()
		return false
	} else {
		r.session = c
		r.available = true
		LogInfo("connect to addr:%s msgque:%d sucess", r.address, r.id)
		if r.handler.OnConnectComplete(r, true) {
			atomic.CompareAndSwapInt32(&r.connecting, 1, 0)
			Go(func() { r.read() })
			Go(func() { r.write() })
			return true
		} else {
			atomic.CompareAndSwapInt32(&r.connecting, 1, 0)
			r.Stop()
			return false
		}
	}
}

func newKcpConn(network, addr string, conn net.Conn, msgtyp MsgType, handler IMsgHandler, parser *Parser, user interface{}) *kcpMsgQue {
	msgque := kcpMsgQue{
		msgQue: msgQue{
			id:            atomic.AddUint32(&msgqueId, 1),
			cwrite:        make(chan *Message, 64),
			msgTyp:        msgtyp,
			handler:       handler,
			timeout:       DefMsgQueTimeout,
			connTyp:       ConnTypeConn,
			parserFactory: parser,
			lastTick:      Timestamp,
			user:          user,
		},
		network: network,
		address: addr,
	}
	if parser != nil {
		msgque.parser = parser.Get()
	}
	msgqueMapSync.Lock()
	msgqueMap[msgque.id] = &msgque
	msgqueMapSync.Unlock()
	LogInfo("new msgque:%d remote addr:%s:%s", msgque.id, network, addr)
	return &msgque
}

func newKcpAccept(session *kcp.UDPSession, msgtyp MsgType, handler IMsgHandler, parser *Parser) *kcpMsgQue {
	msgque := kcpMsgQue{
		msgQue: msgQue{
			id:        atomic.AddUint32(&msgqueId, 1),
			cwrite:    make(chan *Message, cwrite_chan_len),
			msgTyp:    msgtyp,
			handler:   handler,
			available: true,
			timeout:   10,
			connTyp:   ConnTypeAccept,
			//gmsgId:        gmsgId,
			parserFactory: parser,
			lastTick:      Timestamp,
		},
		session: session,
		cread:   make(chan []byte, cread_chan_len),
		addr:    session.LocalAddr(),
	}

	if parser != nil {
		msgque.parser = parser.Get()
	}
	msgqueMapSync.Lock()
	msgqueMap[msgque.id] = &msgque
	msgqueMapSync.Unlock()

	LogInfo("new msgque id:%d from addr:%s", msgque.id, session.RemoteAddr().String())
	return &msgque
}

func newKcpListen(lis *kcp.Listener, msgtyp MsgType, handler IMsgHandler, parser *Parser, addr string) *kcpMsgQue {
	msgque := kcpMsgQue{
		msgQue: msgQue{
			id:            atomic.AddUint32(&msgqueId, 1),
			msgTyp:        msgtyp,
			handler:       handler,
			available:     true,
			parserFactory: parser,
			connTyp:       ConnTypeListen,
		},
		lis:  lis,
		addr: lis.Addr(),
	}
	lis.SetReadBuffer(1 << 24)
	lis.SetWriteBuffer(1 << 24)
	msgqueMapSync.Lock()
	msgqueMap[msgque.id] = &msgque
	msgqueMapSync.Unlock()
	LogInfo("new kcp listen id:%d addr:%s", msgque.id, addr)
	return &msgque
}
