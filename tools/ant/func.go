package antnet

import (
	"core/pb"
	"core/tools/core"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

func AddStopCheck(cs string) uint64 {
	id := atomic.AddUint64(&stopCheckIndex, 1)
	if id == 0 {
		id = atomic.AddUint64(&stopCheckIndex, 1)
	}
	stopCheckMap.Lock()
	stopCheckMap.M[id] = cs
	stopCheckMap.Unlock()
	return id
}

func RemoveStopCheck(id uint64) {
	stopCheckMap.Lock()
	delete(stopCheckMap.M, id)
	stopCheckMap.Unlock()
}

func AtExit(fun func()) {
	id := atomic.AddUint32(&atexitId, 1)
	if id == 0 {
		id = atomic.AddUint32(&atexitId, 1)
	}

	atexitMapSync.Lock()
	atexitMap[id] = fun
	atexitMapSync.Unlock()
}

func Stop() {
	if !atomic.CompareAndSwapInt32(&stop, 0, 1) {
		return
	}
	close(stopChanForGo)
	for sc := 0; !waitAll.TryWait(); sc++ {
		Sleep(1)
		if sc >= 8000 {
			LogError("Server Stop Timeout")
			infos := goInfo.All()
			LogInfo("go not yet closed len: %d", len(infos))
			for id, info := range infos {
				LogInfo("go not yet closed: id: %d, info: %s", id, info)
			}
			stopCheckMap.Lock()
			for _, v := range stopCheckMap.M {
				LogError("Server Stop Timeout:%v", v)
			}
			stopCheckMap.Unlock()
			sc = 0
		}
	}

	LogInfo("Server Stop")
	close(stopChanForSys)
}

func IsStop() bool {
	return stop == 1
}

func IsRuning() bool {
	return stop == 0
}

func CmdAct(cmd, act uint8) int {
	return int(cmd)<<8 + int(act)
}

//func CMDACT(cmd pb.CMD, act pb.ACT) int {
//	return int(cmd)<<8 + int(act)
//}

func Tag(cmd, act uint8) uint16 {
	return uint16(cmd)<<8 | uint16(act)
}

func MD5Str(s string) string {
	return MD5Bytes([]byte(s))
}

func MD5Bytes(s []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(s)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func MD5File(path string) string {
	data, err := ReadFile(path)
	if err != nil {
		LogError("calc md5 failed path:%v", path)
		return ""
	}
	return MD5Bytes(data)
}

func WaitForSystemExit(atexit ...func()) {
	Statistics.startTime = time.Now()
	signal.Notify(stopChanForSys, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-stopChanForSys:
		for _, v := range atexit {
			v()
		}
		//LogInfo("555555555555555555")
		Stop()
		//LogInfo("6666666666666666666")
	}
}

func Daemon(skip ...string) {
	if os.Getppid() != 1 {
		filePath, _ := filepath.Abs(os.Args[0])
		newCmd := []string{}
		for _, v := range os.Args {
			add := true
			for _, s := range skip {
				if strings.Contains(v, s) {
					add = false
					break
				}
			}
			if add {
				newCmd = append(newCmd, v)
			}
		}
		cmd := exec.Command(filePath)
		cmd.Args = newCmd
		cmd.Start()
	}
}

func AtoERROR(str string) pb.ERROR {
	i, err := strconv.Atoi(str)
	if err != nil {
		LogError("AtoERROR error: %v", str)
	}
	return pb.ERROR(i)
}

func Atoi(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func Atoi32(str string) int32 {
	return int32(Atoi(str))
}

func Atoi64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {

		LogError("str to int64 failed err:%v", string(debug.Stack()))
		return 0
	}
	return i
}

func Atof(str string) float32 {
	i, err := strconv.ParseFloat(str, 32)
	if err != nil {
		LogError("str to int64 failed err:%v", err)
		return 0
	}
	return float32(i)
}

func Atof64(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		LogError("str to int64 failed err:%v", err)
		return 0
	}
	return i
}

func Itoa(num interface{}) string {
	switch n := num.(type) {
	case int8:
		return strconv.FormatInt(int64(n), 10)
	case int16:
		return strconv.FormatInt(int64(n), 10)
	case int32:
		return strconv.FormatInt(int64(n), 10)
	case int:
		return strconv.FormatInt(int64(n), 10)
	case int64:
		return strconv.FormatInt(int64(n), 10)
	case uint8:
		return strconv.FormatUint(uint64(n), 10)
	case uint16:
		return strconv.FormatUint(uint64(n), 10)
	case uint32:
		return strconv.FormatUint(uint64(n), 10)
	case uint:
		return strconv.FormatUint(uint64(n), 10)
	case uint64:
		return strconv.FormatUint(uint64(n), 10)
	}
	return ""
}

// 用此函数必然会打印堆栈信息以及错误信息,并且会加入统计
func Try(fun func(), handler func(stack string, e error)) {
	defer core.Exception(handler, func(stack string, e error) { //antnet Try
		LogError(stack)
		Statistics.AddPanic(stack)
	})
	fun()
}

func ParseBaseKind(kind reflect.Kind, data string) (interface{}, error) {
	switch kind {
	case reflect.String:
		return data, nil
	case reflect.Bool:
		v := data == "1" || data == "true"
		return v, nil
	case reflect.Int:
		x, err := strconv.ParseInt(data, 0, 64)
		return int(x), err
	case reflect.Int8:
		x, err := strconv.ParseInt(data, 0, 8)
		return int8(x), err
	case reflect.Int16:
		x, err := strconv.ParseInt(data, 0, 16)
		return int16(x), err
	case reflect.Int32:
		x, err := strconv.ParseInt(data, 0, 32)
		return int32(x), err
	case reflect.Int64:
		x, err := strconv.ParseInt(data, 0, 64)
		return int64(x), err
	case reflect.Float32:
		x, err := strconv.ParseFloat(data, 32)
		return float32(x), err
	case reflect.Float64:
		x, err := strconv.ParseFloat(data, 64)
		return float64(x), err
	case reflect.Uint:
		x, err := strconv.ParseUint(data, 10, 64)
		return uint(x), err
	case reflect.Uint8:
		x, err := strconv.ParseUint(data, 10, 8)
		return uint8(x), err
	case reflect.Uint16:
		x, err := strconv.ParseUint(data, 10, 16)
		return uint16(x), err
	case reflect.Uint32:
		x, err := strconv.ParseUint(data, 10, 32)
		return uint32(x), err
	case reflect.Uint64:
		x, err := strconv.ParseUint(data, 10, 64)
		return uint64(x), err
	default:
		LogError("parse failed type not found type:%v data:%v", kind, data)
		return nil, errors.New("type not found")
	}
}
