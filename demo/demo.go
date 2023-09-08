package main

type Option interface {
	apply(*core)
}
type OptionFunc func(*core)

func (f OptionFunc) apply(core *core) {
	f(core)
}

type lv uint8

func (l lv) GetLv() int32 {
	return int32(l)
}

type LvEnable interface {
	GetLv() int32
}

type core struct {
	lv  int32
	out []string
}

func Level(ab LvEnable) Option {
	return OptionFunc(func(core2 *core) {
		core2.lv = ab.GetLv()
	})
}
func newCore(options ...Option) {
	defCore := &core{}
	for _, option := range options {
		option.apply(defCore)
	}
}

func main1() {
	var k = lv(1)
	newCore(Level(k))
}
