package antnet

import (
	"bufio"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var NetMessageCompressThreshold uint32 = 1000 //网络消息压缩阈值

type tcpMsgQue struct {
	msgQue
	conn       net.Conn     //连接
	listener   net.Listener //监听
	network    string
	address    string
	wait       sync.WaitGroup
	connecting int32
}

func (r *tcpMsgQue) GetNetType() NetType {
	return NetTypeTcp
}

func (r *tcpMsgQue) Wait() {
	r.wait.Wait()
}

func (r *tcpMsgQue) Stop() {
	if atomic.CompareAndSwapInt32(&r.stop, 0, 1) {
		Go(func() {
			if r.init {
				r.handler.OnDelMsgQue(r)
				if r.connecting == 1 {
					r.available = false
					return
				}
			}
			r.available = false
			r.baseStop()
			if r.conn != nil {
				r.conn.Close()
			}
			if r.listener != nil {
				r.listener.Close()
			}
		})
	}
}

func (r *tcpMsgQue) IsStop() bool {
	if r.stop == 0 {
		if IsStop() {
			r.Stop()
		}
	}
	return r.stop == 1
}

func (r *tcpMsgQue) LocalAddr() string {
	if r.conn != nil {
		return r.conn.LocalAddr().String()
	} else if r.listener != nil {
		return r.listener.Addr().String()
	}
	return ""
}

func (r *tcpMsgQue) RemoteAddr() string {
	if r.conn != nil {
		return r.conn.RemoteAddr().String()
	}
	return ""
}

func (r *tcpMsgQue) RemoteIP() string {
	addr := strings.Split(r.RemoteAddr(), ":")
	return addr[0]
}

func (r *tcpMsgQue) readMsg() {
	headData := make([]byte, MsgHeadSize)
	var data []byte
	var head *MessageHead

	for !r.IsStop() {
		if head == nil {
			_, err := io.ReadFull(r.conn, headData)
			if err != nil {
				if err != io.EOF {
					LogDebug("msgque:%v recv data err:%v", r.id, err)
				}
				break
			}
			if head = NewMessageHead(headData); head == nil {
				LogError("msgque:%v read msg head failed", r.id)
				break
			}

			//LogInfo("message head=================%d, %d, %d, %d", head.Len, head.Cmd, head.Act, head.CmdAct())
			if head.Len == 0 {
				if !r.processMsg(r, &Message{Head: head}) {
					LogError("msgque:%v process msg cmd:%v act:%v", r.id, head.Cmd, head.Act)
					break
				}
				head = nil
			} else {
				data = make([]byte, head.Len)
			}
		} else {
			_, err := io.ReadFull(r.conn, data)
			if err != nil {
				LogError("msgque:%v recv data err:%v", r.id, err)
				break
			}

			//LogInfo("send message gamer: %v len: %d, cmd: %d, act: %d, msgId: %d, data: %v", r.user, head.Len, head.Cmd, head.Act, head.CmdAct(), data)
			if !r.processMsg(r, &Message{Head: head, Data: data}) {
				LogError("msgque process failed: %v, cmd: %v, act: %v", r.id, head.Cmd, head.Act)
				break
			}

			head = nil
			data = nil
		}
		r.lastTick = Timestamp
	}
}

func (r *tcpMsgQue) writeMsgFast() {
	var m *Message
	var data []byte
	//gm := r.getGMsg(false)
	writeCount := 0
	tick := time.NewTimer(time.Second * time.Duration(r.timeout))
	for !r.IsStop() || m != nil {
		if m == nil {
			select {
			case <-stopChanForGo:
			case m = <-r.cwrite:
				if m != nil {
					data = m.Bytes()
					//LogInfo("BigEndian==============%v", binary.BigEndian.Uint32(data))
					//LogInfo("LittleEndian==============%v", binary.LittleEndian.Uint32(data))
				}
			//case <-gm.c:
			//	if gm.fun == nil || gm.fun(r) {
			//		m = gm.msg
			//		data = m.Bytes()
			//	}
			//	gm = r.getGMsg(true)
			case <-tick.C:
				if r.isTimeout(tick) {
					r.Stop()
				}
			}
		}

		if m == nil {
			continue
		}

		if writeCount < len(data) {
			n, err := r.conn.Write(data[writeCount:])
			if err != nil {
				LogError("msgque write id:%v err:%v", r.id, err)
				break
			}
			writeCount += n
		}

		if writeCount == len(data) {
			writeCount = 0
			m = nil
		}
		r.lastTick = Timestamp
	}
	tick.Stop()
}

func (r *tcpMsgQue) writeMsg() {
	var m *Message
	//head := make([]byte, MsgHeadSize)
	var data []byte
	//gm := r.getGMsg(false)
	writeCount := 0
	tick := time.NewTimer(time.Second * time.Duration(r.timeout))
	defer tick.Stop()
	for !r.IsStop() || m != nil {
		if m == nil {
			select {
			case <-stopChanForGo:
			case m = <-r.cwrite:
				if m != nil {
					// 广播型数据（如gm服务器发过来的广播消息、KickAllGamer等）是共用一个结构，需要走copy逻辑，避免数据在多个协程中并发
					m = m.Copy()
					//LogInfo("send message pre len: %d, cmd: %d, act: %d, msgId: %d, data: %v", m.Head.Len, m.Head.Cmd, m.Head.Act, m.Head.CmdAct(), m.Data)
					preLen := len(m.Data) + MsgHeadSize
					if r.encrypt {
						preLen += 1
						m.AddDataPlaceholder()
						//LogInfo("222222222222222222222222=============: cmd: %v, act: %v, len: %v, data: %v", m.Cmd(), m.Act(), m.Len(), m.Data)
					}

					if m.Head != nil && m.Head.Len >= NetMessageCompressThreshold && (m.Head.Flags&FlagCompress) == 0 {
						LogDebug("Start GZipCompress head id: %v, len:%v [%v %v]", r.id, m.Head.Len, m.Head.Cmd, m.Head.Act)
						m.Head.Flags |= FlagCompress
						m.Data = GZipCompress(m.Data)
						m.Head.Len = uint32(len(m.Data))
						LogDebug("End GZipCompress head id: %v len:%v", r.id, m.Head.Len)
					}

					//LogInfo(":==================%v", m)
					if r.encrypt && m.Head != nil && m.Head.Cmd != 0 && m.Head.Act != 0 {
						//m = m.Copy()
						//LogInfo("Start Encrypt head:%v data:%v", m.Head, m.Data)

						m.Head.Flags |= FlagEncrypt
						r.oseed = r.oseed*cryptA + cryptB
						m.Head.Bcc = CountBCC(m.Data, 0, m.Head.Len)
						m.Data = DefaultNetEncrypt(r.oseed, m.Data, 0, m.Head.Len)
						//LogInfo("End Encrypt seed:%d bcc:%v head:%v data:%v", r.oseed, m.Head.Bcc, m.Head, m.Data)
					}
					//m.Head.FastBytes(head)
					data = m.Bytes()
					if gid, ok := r.user.(int64); ok && gid > 4000 {
						LogDebug("send message gid: %v, cmd: %d, act: %d, len: %v, %v", r.user, m.Head.Cmd, m.Head.Act, preLen, len(data))
					}
				}
			//case <-gm.c:
			//	if gm.fun == nil || gm.fun(r) {
			//		m = gm.msg
			//		m.Head.FastBytes(head)
			//	}
			//	gm = r.getGMsg(true)
			case <-tick.C:
				if r.isTimeout(tick) {
					r.Stop()
				}
			}
		}

		if m == nil {
			continue
		}

		if writeCount < len(data) {
			n, err := r.conn.Write(data[writeCount:])
			if err != nil {
				LogError("msgque write id:%v err:%v", r.id, err)
				break
			}
			writeCount += n
		}

		if writeCount == len(data) {
			writeCount = 0
			m = nil
		}

		//if writeCount < MsgHeadSize {
		//	n, err := r.conn.Write(head[writeCount:])
		//	if err != nil {
		//		LogError("msgque write id:%v err:%v", r.id, err)
		//		break
		//	}
		//	writeCount += n
		//}
		////LogInfo("33333333333333: %v", m.Data)
		//if writeCount >= MsgHeadSize && m.Data != nil {
		//	n, err := r.conn.Write(m.Data[writeCount-MsgHeadSize : int(m.Head.Len)])
		//	if err != nil {
		//		LogError("msgque write id:%v err:%v", r.id, err)
		//		break
		//	}
		//	writeCount += n
		//}
		////LogInfo("send message len: %d, cmd: %d, act: %d, msgId: %d, data: %v", m.Head.Len, m.Head.Cmd, m.Head.Act, m.Head.CmdAct(), m.Data)
		//if writeCount == int(m.Head.Len)+MsgHeadSize {
		//	/*gid, ok := r.GetUser().(int64)
		//	if ok {
		//		LogInfo("gid:%v,msg sendend cmd:%v,act:%v", gid, m.Head.Cmd, m.Head.Act)
		//	} else {
		//		LogInfo("msg sendend cmd:%v,act:%v", m.Head.Cmd, m.Head.Act)
		//	}*/
		//	writeCount = 0
		//	m = nil
		//}
		r.lastTick = Timestamp
	}
}

func (r *tcpMsgQue) readCmd() {
	reader := bufio.NewReader(r.conn)
	for !r.IsStop() {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		if !r.processMsg(r, &Message{Data: data}) {
			break
		}
		r.lastTick = Timestamp
	}
}

func (r *tcpMsgQue) writeCmd() {
	var m *Message
	//gm := r.getGMsg(false)
	writeCount := 0
	tick := time.NewTimer(time.Second * time.Duration(r.timeout))
	for !r.IsStop() || m != nil {
		if m == nil {
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
					r.Stop()
				}
			}
		}

		if m == nil {
			continue
		}
		n, err := r.conn.Write(m.Data[writeCount:])
		if err != nil {
			LogError("msgque write id:%v err:%v", r.id, err)
			break
		}
		writeCount += n
		if writeCount == len(m.Data) {
			writeCount = 0
			m = nil
		}
		r.lastTick = Timestamp
	}
	tick.Stop()
}

func (r *tcpMsgQue) read() {
	defer func() {
		r.wait.Done()
		if err := recover(); err != nil {
			LogError("msgque read panic id:%v err:%v", r.id, err)
			LogStack()
		}
		r.Stop()
	}()

	r.wait.Add(1)
	if r.msgTyp == MsgTypeCmd {
		r.readCmd()
	} else {
		r.readMsg()
	}
}

func (r *tcpMsgQue) write() {
	defer func() {
		r.wait.Done()
		if err := recover(); err != nil {
			LogError("msgque write panic id:%v err:%v", r.id, err)
			LogStack()
		}
		if r.conn != nil {
			r.conn.Close()
		}
		r.Stop()
	}()
	r.wait.Add(1)
	if r.msgTyp == MsgTypeCmd {
		r.writeCmd()
	} else {
		if r.sendFast {
			//LogInfo("1111111111111111")
			r.writeMsgFast()
		} else {
			//LogInfo("2222222222222222")
			r.writeMsg()
		}
	}
}

func (r *tcpMsgQue) listen() {
	c := make(chan struct{})
	Go2(func(cstop chan struct{}) {
		select {
		case <-cstop:
		case <-c:
		}
		r.listener.Close()
	})
	for !r.IsStop() {
		c, err := r.listener.Accept()
		if err != nil {
			if stop == 0 && r.stop == 0 {
				LogError("accept failed msgque:%v err:%v", r.id, err)
			}
			break
		} else {
			Go(func() {
				msgque := newTcpAccept(c, r.msgTyp, r.handler, r.parserFactory)
				msgque.SetEncrypt(r.GetEncrypt())
				if r.handler.OnNewMsgQue(msgque) {
					msgque.init = true
					msgque.available = true
					Go(func() {
						LogInfo("process read for msgque:%d", msgque.id)
						msgque.read()
						LogInfo("process read end for msgque:%d", msgque.id)
					})
					Go(func() {
						LogInfo("process write for msgque:%d", msgque.id)
						msgque.write()
						LogInfo("process write end for msgque:%d", msgque.id)
					})
				} else {
					msgque.Stop()
				}
			})
		}
	}

	close(c)
	r.Stop()
}

func (r *tcpMsgQue) Connect() bool {
	LogInfo("connect to addr:%s msgque:%d start", r.address, r.id)
	c, err := net.DialTimeout(r.network, r.address, time.Second*3)
	if err != nil {
		LogError("connect to addr:%s msgque:%d err:%v", r.address, r.id, err)
		r.handler.OnConnectComplete(r, false)
		atomic.CompareAndSwapInt32(&r.connecting, 1, 0)
		r.Stop()
		return false
	} else {
		r.conn = c
		r.available = true
		LogInfo("connect to addr:%s msgque:%d sucess", r.address, r.id)
		if r.handler.OnConnectComplete(r, true) {
			atomic.CompareAndSwapInt32(&r.connecting, 1, 0)
			if r.encrypt {
				r.waitRecvSeed = make(chan struct{}, 0)
			}
			Go(func() { r.read() })
			Go(func() { r.write() })

			if r.encrypt {
				select {
				case <-r.waitRecvSeed:
					return true
				case <-time.After(time.Second * 3):
					LogError("connect to addr:%s msgque:%d waitRecvSeed timeout", r.address, r.id)
					return false
				}
			}
			return true
		} else {
			atomic.CompareAndSwapInt32(&r.connecting, 1, 0)
			r.Stop()
			LogError("connect to addr:%s msgque:%d OnConnectComplete failed", r.address, r.id)
			return false
		}
	}
}

func newTcpConn(network, addr string, conn net.Conn, msgtyp MsgType, handler IMsgHandler, parser *Parser, user interface{}) *tcpMsgQue {
	msgque := tcpMsgQue{
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
		conn:    conn,
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

func newTcpAccept(conn net.Conn, msgtyp MsgType, handler IMsgHandler, parser *Parser) *tcpMsgQue {
	msgque := tcpMsgQue{
		msgQue: msgQue{
			id:      atomic.AddUint32(&msgqueId, 1),
			cwrite:  make(chan *Message, 64),
			msgTyp:  msgtyp,
			handler: handler,
			timeout: DefMsgQueTimeout,
			connTyp: ConnTypeAccept,
			//gmsgId:        gmsgId,
			lastTick:      Timestamp,
			parserFactory: parser,
		},
		conn: conn,
		//network: conn.RemoteAddr().Network(),
		//address: conn.RemoteAddr().String(),
	}
	if parser != nil {
		msgque.parser = parser.Get()
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		LogError("TCPConn error")
	}
	tcpConn.SetNoDelay(true)
	msgqueMapSync.Lock()
	msgqueMap[msgque.id] = &msgque
	msgqueMapSync.Unlock()
	LogInfo("new msgque id:%d from addr:%s", msgque.id, conn.RemoteAddr().String())
	return &msgque
}

func newTcpListen(listener net.Listener, msgtyp MsgType, handler IMsgHandler, parser *Parser, addr string) *tcpMsgQue {
	msgque := tcpMsgQue{
		msgQue: msgQue{
			id:            atomic.AddUint32(&msgqueId, 1),
			msgTyp:        msgtyp,
			handler:       handler,
			parserFactory: parser,
			connTyp:       ConnTypeListen,
		},
		listener: listener,
	}

	msgqueMapSync.Lock()
	msgqueMap[msgque.id] = &msgque
	msgqueMapSync.Unlock()
	LogInfo("new tcp listen id:%d addr:%s", msgque.id, addr)
	return &msgque
}
