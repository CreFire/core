// ----------------
// Func  : cluster
// Author: jyb
// Date  : 2021/07/29
// Note  :
// ----------------
package antnet

import (
	"demo/tools/core"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var formatter = core.NewFormatter('@', '<', '>', false)
var scripts = map[string]*Script{}

func GetScript(hash string) *Script {
	return scripts[hash]
}

type Script struct {
	script *redis.Script
	lua    string
}

func (m *Script) String() string {
	return "\n-------------------redis pb--------------------\n" + m.lua + "\n---------------------------------------------------"
}

func NewScript(lua string, kws core.KwArgs, args ...interface{}) *Script {
	str := strings.TrimSpace(lua)
	src := fmt.Sprintf(formatter.Format(str, kws), args...)
	s := &Script{script: redis.NewScript(src), lua: src}
	_, ok := scripts[s.script.Hash()]
	if ok {
		panic(errors.Errorf("new pb exists: %v", s))
	}
	scripts[s.script.Hash()] = s
	return s
}

type ScriptInt64 struct {
	*Script
}

func NewScriptInt64(lua string, kws core.KwArgs, args ...interface{}) *ScriptInt64 {
	return &ScriptInt64{Script: NewScript(lua, kws, args...)}
}

type ScriptString struct {
	*Script
}

func NewScriptString(lua string, kws core.KwArgs, args ...interface{}) *ScriptString {
	return &ScriptString{Script: NewScript(lua, kws, args...)}
}

type ScriptSlice struct {
	*Script
}

func NewScriptSlice(lua string, kws core.KwArgs, args ...interface{}) *ScriptSlice {
	return &ScriptSlice{Script: NewScript(lua, kws, args...)}
}

type ScriptSliceInt64 struct {
	*ScriptSlice
}

func NewScriptSliceInt64(lua string, kws core.KwArgs, args ...interface{}) *ScriptSliceInt64 {
	return &ScriptSliceInt64{ScriptSlice: NewScriptSlice(lua, kws, args...)}
}

type ScriptSliceFloat64 struct {
	*ScriptSlice
}

func NewScriptSliceFloat64(lua string, kws core.KwArgs, args ...interface{}) *ScriptSliceFloat64 {
	return &ScriptSliceFloat64{ScriptSlice: NewScriptSlice(lua, kws, args...)}
}

type ScriptSliceString struct {
	*ScriptSlice
}

func NewScriptSliceString(lua string, kws core.KwArgs, args ...interface{}) *ScriptSliceString {
	return &ScriptSliceString{ScriptSlice: NewScriptSlice(lua, kws, args...)}
}

type ClusterConfig struct {
	Addrs              []string
	Pwd                string
	PoolSize           int
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	PoolTimeout        time.Duration
}

type Cluster struct {
	*redis.ClusterClient
	pubSub   *redis.PubSub
	conf     *ClusterConfig
	channels []string
	closed   int32
	wg       sync.WaitGroup
	//ParallelController *basal.ParallelController
	isKeeWiDB bool
}

func (c *Cluster) Init() error {
	res, err := c.Info().Result()
	if err != nil {
		return err
	}
	idx := strings.Index(res, "# Tx")
	if idx != -1 {
		c.isKeeWiDB = true
	}
	return nil
}

func (c *Cluster) IsKeeWiDB() bool {
	return c.isKeeWiDB
}

func (c *Cluster) Close() {
	if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		c.ClusterClient.Close()
		c.wg.Wait()
	}
}

func (c *Cluster) PrintInfo() {
	LogInfo("Cluster Info: %v, %v", c.conf, c.channels)
}

func (c *Cluster) Unsubscribe(channel ...string) {
	if c.pubSub != nil {
		c.pubSub.Unsubscribe(channel...)
	}
}

func (c *Cluster) SubAdd(channel ...string) error {
	if c.pubSub == nil {
		LogError("Cluster SubAdd pubSub is nil: %v", channel)
		return nil
	}
	return c.pubSub.Subscribe(channel...)
}

func (c *Cluster) Sub(handler func(channel, data string), channels ...string) {
	if len(channels) > 0 {
		if c.pubSub != nil {
			c.pubSub.Close()
		}
		c.channels = channels
		c.pubSub = c.ClusterClient.Subscribe(channels...)
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			for atomic.LoadInt32(&c.closed) == 0 {
				msg, err := c.pubSub.ReceiveMessage()
				if err == nil {
					Try(func() {
						handler(msg.Channel, msg.Payload)
					}, func(stack string, e error) {
						LogError("Cluster ReceiveMessage err: %v", e)
					})
				} else {
					LogError("Cluster ReceiveMessage err2: %v", err)
					break
				}
			}
			LogInfo("Cluster ReceiveMessage Stopped")
		}()
	} else {
		LogError("Cluster Sub Channels len == 0")
	}
}

func (c *Cluster) Script(s *Script, keys []string, args ...interface{}) (interface{}, error) {
	return s.script.Run(c, keys, args...).Result()
}

func (c *Cluster) ScriptInt64(s *ScriptInt64, keys []string, args ...interface{}) (int64, error) {
	result, err := c.Script(s.Script, keys, args...)
	if err != nil {
		return 0, err
	}
	if res, ok := result.(int64); ok {
		return res, nil
	}
	return 0, errors.Errorf("redis result type is %v", reflect.TypeOf(result))
}

var ErrResultNotString = errors.New("redis result not is string")

func (c *Cluster) ScriptString(s *ScriptString, keys []string, args ...interface{}) (string, error) {
	result, err := c.Script(s.Script, keys, args...)
	if err != nil {
		return "", err
	}
	if res, ok := result.(string); ok {
		return res, nil
	}
	return "", ErrResultNotString
}

var ErrResultNotSlice = errors.New("redis result not is []interface")

func (c *Cluster) ScriptSlice(s *ScriptSlice, keys []string, args ...interface{}) ([]interface{}, error) {
	result, err := c.Script(s.Script, keys, args...)
	if err != nil {
		return nil, err
	}
	if res, ok := result.([]interface{}); ok {
		return res, nil
	}
	return nil, ErrResultNotSlice
}

var ErrResultNotSliceInt64 = errors.New("redis result not is []int64")

func (c *Cluster) ScriptSliceInt64(s *ScriptSliceInt64, keys []string, args ...interface{}) ([]int64, error) {
	result, err := c.ScriptSlice(s.ScriptSlice, keys, args...)
	if err != nil {
		return nil, err
	}
	res := make([]int64, len(result))
	for index, a := range result {
		if v, ok := a.(int64); ok {
			res[index] = v
		} else {
			return nil, ErrResultNotSliceInt64
		}
	}
	return res, nil
}

var ErrResultNotSliceFloat64 = errors.New("redis result not is []float64")

func (c *Cluster) ScriptSliceFloat64(s *ScriptSliceFloat64, keys []string, args ...interface{}) ([]float64, error) {
	result, err := c.ScriptSlice(s.ScriptSlice, keys, args...)
	if err != nil {
		return nil, err
	}
	res := make([]float64, len(result))
	for index, a := range result {
		if v, ok := a.(float64); ok {
			res[index] = v
		} else {
			return nil, ErrResultNotSliceFloat64
		}
	}
	return res, nil
}

var ErrResultNotSliceString = errors.New("redis result not is []string")

func (c *Cluster) ScriptSliceString(s *ScriptSliceString, keys []string, args ...interface{}) ([]string, error) {
	result, err := c.ScriptSlice(s.ScriptSlice, keys, args...)
	if err != nil {
		return nil, err
	}
	res := make([]string, len(result))
	for index, a := range result {
		if v, ok := a.(string); ok {
			res[index] = v
		} else {
			return nil, ErrResultNotSliceString
		}
	}
	return res, nil
}

func NewCluster(conf *ClusterConfig) *Cluster {
	c := &Cluster{conf: conf}
	opt := &redis.ClusterOptions{
		Addrs:              conf.Addrs,
		Password:           conf.Pwd,
		PoolSize:           conf.PoolSize,
		IdleTimeout:        conf.IdleTimeout,
		IdleCheckFrequency: conf.IdleCheckFrequency,
		PoolTimeout:        conf.PoolTimeout,
	}
	c.ClusterClient = redis.NewClusterClient(opt)
	//c.ParallelController = basal.NewParallelController(conf.PoolSize + 10)
	return c
}

//func (c *Cluster) Pipeline() redis.Pipeliner {
//	return NewPipeline(c.ClusterClient, 0)
//}
