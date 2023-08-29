package antnet

import (
	"core/tools/core"
	"sync"
	"sync/atomic"
)

func (m *goCreatedInfo) Set(id uint32, info string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.infos[id] = info
}

func (m *goCreatedInfo) Get(id uint32) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.infos[id]
	return v, ok
}

func (m *goCreatedInfo) Del(id uint32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.infos, id)
}

func (m *goCreatedInfo) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.infos)
}

func (m *goCreatedInfo) All() map[uint32]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	infos := make(map[uint32]string, len(m.infos))
	for id, info := range m.infos {
		infos[id] = info
	}
	return infos
}

var goInfo = &goCreatedInfo{infos: make(map[uint32]string)}

func Go(fn func()) {
	waitAll.Add(1)
	var debugStr string
	id := atomic.AddUint32(&goid, 1)
	c := atomic.AddInt32(&gocount, 1)
	if DefLog.Level() <= LogLevelDebug {
		debugStr = LogSimpleStack()
		LogDebug("goroutine start id:%d count:%d from:%s", id, c, debugStr)
	}

	name, file, line := core.CallerInFunc(2)
	goInfo.Set(id, Sprintf("%s, %s:%d", name, file, line))
	//LogInfo("go created by %v, %v:%v", name, file, line)
	go func() {
		Try(fn, nil)
		waitAll.Done()
		c = atomic.AddInt32(&gocount, -1)
		if DefLog.Level() <= LogLevelDebug {
			LogDebug("goroutine end id:%d count:%d from:%s", id, c, debugStr)
		}
		goInfo.Del(id)
		//LogInfo("go closed by %v, %v:%v", name, file, line)
	}()
}

type goCreatedInfo struct {
	mu    sync.RWMutex
	infos map[uint32]string
}

func Go2(fn func(cstop chan struct{})) bool {
	if IsStop() {
		return false
	}
	waitAll.Add(1)
	var debugStr string
	id := atomic.AddUint32(&goid, 1)
	c := atomic.AddInt32(&gocount, 1)
	if DefLog.Level() <= LogLevelDebug {
		debugStr = LogSimpleStack()
		LogDebug("goroutine start id:%d count:%d from:%s", id, c, debugStr)
	}
	name, file, line := core.CallerInFunc(2)
	goInfo.Set(id, Sprintf("%s, %s:%d", name, file, line))
	//LogInfo("go2 created by %v, %v:%v", name, file, line)
	go func() {
		Try(func() { fn(stopChanForGo) }, nil)
		waitAll.Done()
		c = atomic.AddInt32(&gocount, -1)
		if DefLog.Level() <= LogLevelDebug {
			LogDebug("goroutine end id:%d count:%d from:%s", id, c, debugStr)
		}
		goInfo.Del(id)
		//LogInfo("go2 closed by %v, %v:%v", name, file, line)
	}()
	return true
}

func GoArgs(fn func(...interface{}), args ...interface{}) {
	waitAll.Add(1)
	var debugStr string
	id := atomic.AddUint32(&goid, 1)
	c := atomic.AddInt32(&gocount, 1)
	if DefLog.Level() <= LogLevelDebug {
		debugStr = LogSimpleStack()
		LogDebug("goroutine start id:%d count:%d from:%s", id, c, debugStr)
	}

	go func() {
		Try(func() { fn(args...) }, nil)

		waitAll.Done()
		c = atomic.AddInt32(&gocount, -1)
		if DefLog.Level() <= LogLevelDebug {
			LogDebug("goroutine end id:%d count:%d from:%s", id, c, debugStr)
		}
	}()
}

func goForRedis(fn func()) {
	waitAllForRedis.Add(1)
	var debugStr string
	id := atomic.AddUint32(&goid, 1)
	c := atomic.AddInt32(&gocount, 1)
	if DefLog.Level() <= LogLevelDebug {
		debugStr = LogSimpleStack()
		LogDebug("goroutine start id:%d count:%d from:%s", id, c, debugStr)
	}
	go func() {
		Try(fn, nil)
		waitAllForRedis.Done()
		c = atomic.AddInt32(&gocount, -1)

		if DefLog.Level() <= LogLevelDebug {
			LogDebug("goroutine end id:%d count:%d from:%s", id, c, debugStr)
		}
	}()
}

func goForLog(fn func(cstop chan struct{})) bool {
	if IsStop() {
		return false
	}
	waitAllForLog.Add(1)

	go func() {
		fn(stopChanForLog)
		waitAllForLog.Done()
	}()
	return true
}
