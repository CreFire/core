package antnet

import (
	"core/tools/core"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Addr     string
	Passwd   string
	PoolSize int
}

type Redis struct {
	*redis.Client
	pubsub  *redis.PubSub
	conf    *RedisConfig
	manager *RedisManager
}

func (r *Redis) ScriptStr(cmd int, keys []string, args ...interface{}) (string, error) {
	data, err := r.Script(cmd, keys, args...)
	if err != nil {
		LogError("redis pb failed err:%v", err)
		return "", ErrDBErr
	}
	errcode, ok := data.(int64)
	if ok {
		return "", GetError(uint16(errcode))
	}

	str, ok := data.(string)
	if !ok {
		return "", ErrDBDataType
	}

	return str, nil
}

func (r *Redis) ScriptStrArray(cmd int, keys []string, args ...interface{}) ([]string, error) {
	data, err := r.Script(cmd, keys, args...)
	if err != nil {
		LogError("redis pb failed err:%v", err)
		return nil, ErrDBErr
	}
	errcode, ok := data.(int64)
	if ok {
		return nil, GetError(uint16(errcode))
	}

	iArray, ok := data.([]interface{})
	if !ok {
		return nil, ErrDBDataType
	}

	strArray := []string{}
	for _, v := range iArray {
		if str, ok := v.(string); ok {
			strArray = append(strArray, str)
		} else {
			return nil, ErrDBDataType
		}
	}

	return strArray, nil
}

func (r *Redis) ScriptInt64(cmd int, keys []string, args ...interface{}) (int64, error) {
	data, err := r.Script(cmd, keys, args...)
	if err != nil {
		LogError("redis pb failed err:%v", err)
		return 0, ErrDBErr
	}
	code, ok := data.(int64)
	if ok {
		return code, nil
	}
	return 0, ErrDBDataType
}

func (r *Redis) Script(cmd int, keys []string, args ...interface{}) (interface{}, error) {
	data, _ := scriptHashMap.Load(cmd)
	hash, _ := data.(string)
	// 腾讯云redis必须至少填一个key
	//keys = append(keys, "bug{tag}")
	re, err := r.EvalSha(hash, keys, args...).Result()
	if err != nil {
		script, ok := scriptMap[cmd]
		if !ok {
			LogError("redis pb error cmd not found cmd:%v", cmd)
			return nil, ErrDBErr
		}

		if strings.HasPrefix(err.Error(), "NOSCRIPT ") {
			LogInfo("try reload redis pb %v", scriptCommitMap[cmd])
			hash, err = r.ScriptLoad(script).Result()
			if err != nil {
				LogError("redis pb load cmd:%v errstr:%s", scriptCommitMap[cmd], err)
				return nil, ErrDBErr
			}
			scriptHashMap.Store(cmd, hash)
			re, err = r.EvalSha(hash, keys, args...).Result()
			if err == nil {
				return re, nil
			}
		}
		LogError("redis pb error cmd:%v err:%s", scriptCommitMap[cmd], err)
		return nil, ErrDBErr
	}

	return re, nil
}

type RedisManager struct {
	dbs      map[int]*Redis
	subMap   map[string]*Redis
	channels []string
	fun      func(channel, data string)
	lock     sync.RWMutex
}

func (r *RedisManager) GetByRid(rid int) *Redis {
	r.lock.RLock()
	defer r.lock.RUnlock()
	db, ok := r.dbs[rid]
	if !ok {
		LogError("GetByRid:%v not find", rid)
		panic(ErrDBErr)
	}
	return db
}

func (r *RedisManager) CheckByRid(rid int) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	_, ok := r.dbs[rid]
	if !ok {
		return false
	}
	return true
}

func (r *RedisManager) GetGlobal() *Redis {
	return r.GetByRid(0)
}

func (r *RedisManager) Sub(fun func(channel, data string), channels ...string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.channels = channels
	r.fun = fun
	for _, v := range r.subMap {
		if v.pubsub != nil {
			v.pubsub.Close()
		}
		pubsub := v.Subscribe(channels...)
		v.pubsub = pubsub
		goForRedis(func() {
			for IsRuning() {
				msg, err := pubsub.ReceiveMessage()
				//LogInfo("ReceiveMessage msg:%v pubsub:%v", msg, pubsub)
				if err == nil {
					Go(func() { fun(msg.Channel, msg.Payload) })
				} else if _, ok := err.(net.Error); !ok {
					break
				}
			}
		})
	}
}

func (r *RedisManager) Exist(id int) bool {
	r.lock.Lock()
	_, ok := r.dbs[id]
	r.lock.Unlock()
	return ok
}

func (r *RedisManager) Add(id int, conf *RedisConfig) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if _, ok := r.dbs[id]; ok {
		LogError("redis already have id:%v", id)
		return
	}
	re := &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Passwd,
			PoolSize: conf.PoolSize,
		}),
		conf:    conf,
		manager: r,
	}

	re.WrapProcess(func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			err := oldProcess(cmd)
			if err != nil {
				_, retry := err.(net.Error)
				if !retry {
					retry = err == io.EOF
				}
				if retry {
					err = oldProcess(cmd)
				}
			}
			return err
		}
	})

	if _, ok := r.subMap[conf.Addr]; !ok {
		r.subMap[conf.Addr] = re
		if len(r.channels) > 0 {
			pubsub := re.Subscribe(r.channels...)
			re.pubsub = pubsub
			goForRedis(func() {
				for IsRuning() {
					msg, err := pubsub.ReceiveMessage()
					if err == nil {
						Go(func() { r.fun(msg.Channel, msg.Payload) })
					} else if _, ok := err.(net.Error); !ok {
						break
					}
				}
			})
		}
	}
	r.dbs[id] = re
	LogInfo("connect to redis %v", conf.Addr)
}

func (r *RedisManager) Info() {
	LogInfo("RedisManager dbs:%v submap:%v channels:%v", r.dbs, r.subMap, r.channels)
}

func (r *RedisManager) close() {
	for _, v := range r.dbs {
		if v.pubsub != nil {
			v.pubsub.Close()
		}
		v.Close()
	}
}

var (
	scriptMap       = map[int]string{}
	scriptCommitMap = map[int]string{}
	//scriptHashMap         = map[int]string{}
	scriptIndex   int32 = 0
	scriptHashMap sync.Map
)

var ScriptFormatter = core.NewFormatter('@', '<', '>', false)

func NewRedisScript(commit, str string) int {
	cmd := int(atomic.AddInt32(&scriptIndex, 1))
	scriptMap[cmd] = str
	scriptCommitMap[cmd] = commit
	return cmd
}

var redisManagers []*RedisManager

func NewRedisManager(conf *RedisConfig) *RedisManager {
	redisManager := &RedisManager{
		subMap: map[string]*Redis{},
		dbs:    map[int]*Redis{},
	}

	redisManager.Add(0, conf)
	redisManagers = append(redisManagers, redisManager)
	return redisManager
}

func RedisNil(err error) bool {
	return err == redis.Nil
}

func RedisError(err error) bool {
	if err == redis.Nil {
		return false
	}
	return err != nil
}
