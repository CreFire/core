package core

import (
	"sync"
)

const bufferPoolNumber = 10
const bufferPoolIndexMax = bufferPoolNumber - 1

const minBufferSize = 1 << 10

// buffer size lv
var bufferLv = 1

var GloBufferPool bufferPools

func init() {
	GloBufferPool.init()
}

type bufferPool struct {
	sync.Pool
	size int64
}

func (m *bufferPool) Get() *Buffer {

	return m.Pool.Get().(*Buffer)
}

const uint8Max = (1 << 8) - 1

func (m *bufferPool) init(index int) {
	m.size = int64((index + 1) * minBufferSize * bufferLv)
	m.New = func() any {
		if index > uint8Max {
			index = uint8Max
		}
		return &Buffer{
			index: uint8(index),
			buf:   make([]byte, 0, m.size),
		}
	}
}

type bufferPools [bufferPoolNumber]*bufferPool

func (m *bufferPools) init() {
	for i := 0; i < bufferPoolNumber; i++ {
		m[i] = &bufferPool{}
		m[i].init(i)
	}
}

func (m *bufferPools) New(size int) *Buffer {
	bufferSize := minBufferSize * bufferLv
	index := size / bufferSize
	if size%bufferSize == 0 {
		index--
	}
	//log.InfoF("================New: %v", index)
	if index < 0 {
		return m[0].Get()
	}
	if index < bufferPoolNumber {
		return m[index].Get()
	} else {
		return &Buffer{buf: make([]byte, 0, size), index: uint8(index)}
	}
}
func (m *bufferPools) Free(buf *Buffer) {
	index := buf.index
	if index < bufferPoolNumber {
		m[index].Put(buf)
	} else {
		m[bufferPoolIndexMax].Put(buf)
	}
}
