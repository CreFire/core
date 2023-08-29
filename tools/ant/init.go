package antnet

import (
	"core/pb"
	"core/tools/datetime"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type statistics struct {
	startTime   time.Time
	lastPanic   int64
	panicCount  int32
	PanicStacks chan *pb.InnerNotfiyEmailData //panic堆栈
}

func (m *statistics) AddPanic(stack string) {
	atomic.AddInt32(&m.panicCount, 1)
	atomic.StoreInt64(&m.lastPanic, Timestamp)
	data := &pb.InnerNotfiyEmailData{Data: stack, Dt: datetime.UnixToYmdHMS(datetime.Unix(), datetime.Local())}
	select {
	case m.PanicStacks <- data:
	default:
	}
}

// 开始时间
func (m *statistics) StartTime() time.Time {
	return m.startTime
}

// 添加panic次数
func (m *statistics) AddPanicCount() {
	atomic.AddInt32(&m.panicCount, 1)
	atomic.StoreInt64(&m.lastPanic, Timestamp)
}

// panic数量
func (m *statistics) PanicCount() int32 {
	return atomic.LoadInt32(&m.panicCount)
}

// 最后panic时间
func (m *statistics) LastPanic() int64 {
	return atomic.LoadInt64(&m.lastPanic)
}

// 协程数
func (m *statistics) GoCount() int {
	return int(atomic.LoadInt32(&gocount))
}

// 连接数
func (m *statistics) MsgQueCount() int {
	return len(msgqueMap)
}

// 统计
var Statistics = &statistics{PanicStacks: make(chan *pb.InnerNotfiyEmailData, 4096)}

type WaitGroup struct {
	count int64
}

func (r *WaitGroup) Add(delta int) {
	atomic.AddInt64(&r.count, int64(delta))
}

func (r *WaitGroup) Done() {
	atomic.AddInt64(&r.count, -1)
}

func (r *WaitGroup) Wait() {
	for atomic.LoadInt64(&r.count) > 0 {
		Sleep(1)
	}
}

func (r *WaitGroup) TryWait() bool {
	return atomic.LoadInt64(&r.count) == 0
}

var waitAll = &WaitGroup{} //等待所有goroutine
var waitAllForLog sync.WaitGroup
var waitAllForRedis sync.WaitGroup

// var stopForLog int32 //
var stop int32 //停止标志

var gocount int32 //goroutine数量
var goid uint32
var DefLog *Log    //程序日志
var RecordLog *Log //记录日志

var msgqueId uint32 //消息队列id
var msgqueMapSync sync.Mutex
var msgqueMap = map[uint32]IMsgQue{}

type gMsg struct {
	c   chan struct{}
	msg *Message
	fun func(msgque IMsgQue) bool
}

//var gmsgId uint16
//var gmsgMapSync sync.Mutex
//var gmsgArray = [65536]*gMsg{}

var atexitId uint32
var atexitMapSync sync.Mutex
var atexitMap = map[uint32]func(){}

var stopChanForGo = make(chan struct{})
var stopChanForLog = make(chan struct{})
var stopChanForSys = make(chan os.Signal, 1)

var StartTick int64
var NowTick int64
var Timestamp int64
var StartUnix int64

var UdpServerGoCnt = 64

var stopCheckIndex uint64
var stopCheckMap = struct {
	sync.Mutex
	M map[uint64]string
}{M: map[uint64]string{}}

func GetStopChan() chan struct{} {
	return stopChanForGo
}

func init() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
	DefLog = NewLog(10000)
	DefLog.SetLogger(&ConsoleLogger{true})
	DefLog.SetLevel(LogLevelInfo)
	RecordLog = NewLog(10000)
	timerTick()
	WeekStart = DateToUnix("2018-01-01 00:00:00") //2018/1/1
}
