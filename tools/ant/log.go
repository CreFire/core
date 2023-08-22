package antnet

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

type ILogger interface {
	Write(str string)
}

type ConsoleLogger struct {
	Ln bool
}

func (r *ConsoleLogger) Write(str string) {
	if r.Ln {
		fmt.Println(str)
	} else {
		fmt.Print(str)
	}
}

type OnFileLogFull func(path string)
type OnFileLogTimeout func(path string) int
type FileLogger struct {
	Path      string
	Ln        bool
	Timeout   int //0表示不设置, 单位s
	MaxSize   int //0表示不限制，最大大小
	OnFull    OnFileLogFull
	OnTimeout OnFileLogTimeout

	size     int
	file     *os.File
	filename string
	extname  string
	dirname  string
}

func (r *FileLogger) Write(str string) {
	if r.file == nil {
		return
	}

	newsize := r.size
	if r.Ln {
		newsize += len(str) + 1
	} else {
		newsize += len(str)
	}

	if r.MaxSize > 0 && newsize >= r.MaxSize {
		r.file.Close()
		r.file = nil
		newpath := r.dirname + "/" + r.filename + fmt.Sprintf("_%v", datetime.Now().YmdHMS()) + r.extname
		err := os.Rename(r.Path, newpath)
		if err != nil {
			LogInfo("full rename path:%v newpath:%v err:%v ", r.Path, newpath, err)
		}
		file, err := os.OpenFile(r.Path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err == nil {
			r.file = file
		} else {
			LogInfo("full openfile path:%v err:%v ", r.Path, err)
		}
		r.size = 0
		if r.OnFull != nil {
			r.OnFull(newpath)
		}
	}

	if r.file == nil {
		return
	}

	if r.Ln {
		r.file.WriteString(str)
		r.file.WriteString("\n")
		r.size += len(str) + 1
	} else {
		r.file.WriteString(str)
		r.size += len(str)
	}
}

type LogLevel int

const (
	LogLevelAllOn  LogLevel = iota //开放说有日志
	LogLevelDebug                  //调试信息
	LogLevelInfo                   //资讯讯息
	LogLevelWarn                   //警告状况发生
	LogLevelError                  //一般错误，可能导致功能不正常
	LogLevelFatal                  //严重错误，会导致进程退出
	LogLevelAllOff                 //关闭所有日志
)

type Log struct {
	logger         [32]ILogger
	cwrite         chan string
	ctimeout       chan *FileLogger
	bufsize        int
	stop           int32
	preLoggerCount int32
	loggerCount    int32
	level          LogLevel
	closed         int32
	exitWg         sync.WaitGroup
}

func (r *Log) initFileLogger(f *FileLogger) *FileLogger {
	if f.file == nil {
		f.Path, _ = filepath.Abs(f.Path)
		f.dirname = path.Dir(f.Path)
		f.extname = path.Ext(f.Path)
		f.filename = filepath.Base(f.Path[0 : len(f.Path)-len(f.extname)])
		os.MkdirAll(f.dirname, 0777)
		file, err := os.OpenFile(f.Path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err == nil {
			f.file = file
			info, err := f.file.Stat()
			if err != nil {
				return nil
			}

			f.size = int(info.Size())
			if f.Timeout > 0 {
				SetTimeout(f.Timeout*1000, func(...interface{}) int {
					defer func() { recover() }()
					r.ctimeout <- f
					return 0
				})
			}

			return f
		}
	}
	return nil
}

func (r *Log) run() {
	defer r.exitWg.Done()
	var i int32
	cwrite := r.cwrite
	ctimeout := r.ctimeout
	cwriteOk := true
	ctimeoutOk := true

	for cwriteOk || ctimeoutOk {
		select {
		case s, ok := <-cwrite:
			if ok {
				for i = 0; i < r.loggerCount; i++ {
					r.logger[i].Write(s)
				}
			} else {
				cwriteOk = false
				cwrite = nil
			}
		case c, ok := <-ctimeout:
			if ok {
				c.file.Close()
				c.file = nil
				newpath := c.dirname + "/" + c.filename + fmt.Sprintf("_%v", datetime.Now().YmdHMS()) + c.extname
				err := os.Rename(c.Path, newpath)
				if err != nil {
					LogInfo("timeout rename path:%v newpath:%v err:%v", c.Path, newpath, err)
				}
				file, err := os.OpenFile(c.Path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
				if err == nil {
					c.file = file
				} else {
					LogInfo("timeout open file:%v err:%v", c.Path, err)
				}
				c.size = 0
				if c.OnTimeout != nil {
					nc := c.OnTimeout(newpath)
					if nc > 0 {
						SetTimeout(nc*1000, func(...interface{}) int {
							defer func() { recover() }()
							r.ctimeout <- c
							return 0
						})
					}
				}
			} else {
				ctimeoutOk = false
				ctimeout = nil
			}
		}
	}
}

func (r *Log) start() {
	r.exitWg.Add(1)
	go r.run()
}

//func (r *Log) Stop() {
//	if atomic.CompareAndSwapInt32(&r.stop, 0, 1) {
//		close(r.cwrite)
//		close(r.ctimeout)
//	}
//}

func (r *Log) Close() {
	if atomic.CompareAndSwapInt32(&r.closed, 0, 1) {
		close(r.cwrite)
		close(r.ctimeout)
	}
	r.exitWg.Wait()
}

func (r *Log) SetLogger(logger ILogger) bool {
	if r.preLoggerCount >= 31 {
		return false
	}
	if f, ok := logger.(*FileLogger); ok {
		if r.initFileLogger(f) == nil {
			return false
		}
	}
	r.logger[atomic.AddInt32(&r.preLoggerCount, 1)] = logger
	atomic.AddInt32(&r.loggerCount, 1)
	return true

}
func (r *Log) Level() LogLevel {
	return r.level
}
func (r *Log) SetLevel(level LogLevel) {
	r.level = level
}

//func isLogStop() bool {
//	return stopForLog == 1
//}

//func (r *Log) IsStop() bool {
//	if r.stop == 0 {
//		if isLogStop() {
//			r.Stop()
//		}
//	}
//	return r.stop == 1
//}

func (r *Log) write(levstr string, v ...interface{}) {
	prefix := levstr
	_, file, line, ok := runtime.Caller(3)
	if ok {
		i := strings.LastIndex(file, "/") + 1
		prefix = fmt.Sprintf("[%s][%s][%s:%d]:", levstr, Date(), (string)(([]byte(file))[i:]), line)
	}
	var logText string
	if len(v) > 1 {
		logText = prefix + fmt.Sprintf(v[0].(string), v[1:]...)
	} else {
		logText = prefix + fmt.Sprint(v[0])
	}
	defer func() { recover() }()
	r.cwrite <- logText
}

func (r *Log) Debug(v ...interface{}) {
	if r.level <= LogLevelDebug {
		r.write("D", v...)
	}
}

func (r *Log) Info(v ...interface{}) {
	if r.level <= LogLevelInfo {
		r.write("I", v...)
	}
}

func (r *Log) Warn(v ...interface{}) {
	if r.level <= LogLevelWarn {
		r.write("W", v...)
	}
}

func (r *Log) Error(v ...interface{}) {
	if r.level <= LogLevelError {
		r.write("E", v...)
	}
}

func (r *Log) ErrorF(format string, v ...interface{}) {
	args := make([]interface{}, len(v)+1)
	args[0] = format
	copy(args[1:], v)
	if r.level <= LogLevelError {
		r.write("E", args...)
	}
}

func (r *Log) Fatal(v ...interface{}) {
	if r.level <= LogLevelFatal {
		r.write("FATAL", v...)
	}
}

func (r *Log) Write(v ...interface{}) {
	var logText string
	if len(v) > 1 {
		logText = fmt.Sprintf(v[0].(string), v[1:]...)
	} else if len(v) > 0 {
		logText = fmt.Sprint(v[0])
	}
	defer func() { recover() }()
	r.cwrite <- logText
}

func NewLog(bufsize int, logger ...ILogger) *Log {
	log := &Log{
		bufsize:        bufsize,
		cwrite:         make(chan string, bufsize),
		ctimeout:       make(chan *FileLogger, 32),
		level:          LogLevelAllOn,
		preLoggerCount: -1,
	}
	for _, l := range logger {
		log.SetLogger(l)
	}
	log.start()
	return log
}

func LogInfo(v ...interface{}) {
	DefLog.Info(v...)
}

func LogDebug(v ...interface{}) {
	DefLog.Debug(v...)
}

func LogError(v ...interface{}) {
	DefLog.Error(v...)
}

func LogFatal(v ...interface{}) {
	DefLog.Fatal(v...)
}

func LogWarn(v ...interface{}) {
	DefLog.Warn(v...)
}

func LogStack() {
	buf := make([]byte, 1<<12)
	LogError(string(buf[:runtime.Stack(buf, false)]))
}

func LogSimpleStack() string {
	_, file, line, _ := runtime.Caller(2)
	i := strings.LastIndex(file, "/") + 1
	i = strings.LastIndex((string)(([]byte(file))[:i-1]), "/") + 1

	return Sprintf("%s:%d", (string)(([]byte(file))[i:]), line)
}

func RecordInfoPb(key string, value proto.Message) {
	if js, err := JsonStr(value); err == nil {
		RecordLog.Info(Sprintf("%s:%s", key, js))
	}
}

func RecordRobotInfo(v ...interface{}) {
	RecordLog.Write(v...)
}

type iLogger struct {
}

func (*iLogger) ErrorF(format string, v ...interface{}) {
	args := make([]interface{}, len(v)+1)
	args[0] = format
	copy(args[1:], v)
	DefLog.Error(args...)
}

func (*iLogger) InfoF(format string, v ...interface{}) {
	args := make([]interface{}, len(v)+1)
	args[0] = format
	copy(args[1:], v)
	DefLog.Info(args...)
}

var logger = &iLogger{}

func GetLogger() basal.ILogger {
	return logger
}
