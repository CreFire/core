package core

type Buffer struct {
	index uint8
	buf   []byte
}

func (m *Buffer) AppendByte(b byte) {
	m.buf = append(m.buf, b)
}

func (m *Buffer) AppendBytes(bs ...byte) {
	m.buf = append(m.buf, bs...)
}

func (m *Buffer) AppendString(s string) {
	m.buf = append(m.buf, s...)
}

func (m *Buffer) AppendStrings(ss ...string) {
	for _, s := range ss {
		m.buf = append(m.buf, s...)
	}
}

// Bytes
//
//	@Description: 外部调用了Bytes 后没有回收前不可重用
//	@receiver m
//	@return []byte
func (m *Buffer) Bytes() []byte {
	return m.buf
}

func (m *Buffer) ToString() string {
	return string(m.buf)
}

func (m *Buffer) Clear() {
	m.buf = m.buf[:0]
}

func (m *Buffer) Len() int {
	return len(m.buf)
}

func (m *Buffer) Cap() int {
	return cap(m.buf)
}

func (m *Buffer) Free() {
	GloBufferPool.Free(m)
}

//go:linkname NewBuffer github.com/jingyanbin/core/basal.NewBuffer
func NewBuffer(size int) *Buffer {
	return GloBufferPool.New(size)
}
