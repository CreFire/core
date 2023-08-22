package antnet

import (
	"github.com/go-redis/redis"
	"time"
)

const pipeCmdMax = 20

func NewPipeline(c *redis.ClusterClient, size int) *Pipeline {
	pipe := &Pipeline{c: c}
	if size > 0 {
		pipe.pipes = make([]redis.Pipeliner, 0, size/pipeCmdMax)
	}
	return pipe
}

type Pipeline struct {
	pipes []redis.Pipeliner
	num   uint32
	c     *redis.ClusterClient
}

func (m *Pipeline) getPipe() redis.Pipeliner {
	index := m.num / pipeCmdMax
	m.num += 1
	if index < uint32(len(m.pipes)) {
		return m.pipes[index]
	} else {
		pipe := m.c.Pipeline()
		m.pipes = append(m.pipes, pipe)
		return pipe
	}
}

func (m *Pipeline) Do(args ...interface{}) *redis.Cmd {
	return nil
}

func (m *Pipeline) Process(cmd redis.Cmder) error {
	return nil
}

func (m *Pipeline) Close() (err error) {
	for _, pipe := range m.pipes {
		if e := pipe.Close(); e != nil {
			err = e
		}
	}
	return
}

func (m *Pipeline) Discard() (err error) {
	for _, pipe := range m.pipes {
		if e := pipe.Discard(); e != nil {
			err = e
		}
	}
	return
}

func (m *Pipeline) Exec() ([]redis.Cmder, error) {
	var err error
	var cmdArr = make([]redis.Cmder, 0, m.num)
	for _, pipe := range m.pipes {
		cmds, e := pipe.Exec()
		if e != nil {
			err = e
		}
		cmdArr = append(cmdArr, cmds...)
	}
	return cmdArr, err
}

func (m *Pipeline) PipeNum() int {
	return len(m.pipes)
}

func (m *Pipeline) Auth(password string) *redis.StatusCmd {
	//return m.getPipe().Auth(password)
	return nil
}

func (m *Pipeline) Select(index int) *redis.StatusCmd {
	return m.getPipe().Select(index)
}

func (m *Pipeline) SwapDB(index1, index2 int) *redis.StatusCmd {
	return m.getPipe().SwapDB(index1, index2)
}

func (m *Pipeline) ClientSetName(name string) *redis.BoolCmd {
	return m.getPipe().ClientSetName(name)
}

func (m *Pipeline) Pipeline() redis.Pipeliner {
	return nil
}

func (m *Pipeline) Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return nil, nil
}

func (m *Pipeline) TxPipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return nil, nil
}

func (m *Pipeline) TxPipeline() redis.Pipeliner {
	return nil
}

func (m *Pipeline) Command() *redis.CommandsInfoCmd {
	return m.getPipe().Command()
}

func (m *Pipeline) ClientGetName() *redis.StringCmd {
	return m.getPipe().ClientGetName()
}

func (m *Pipeline) Echo(message interface{}) *redis.StringCmd {
	return m.getPipe().Echo(message)
}

func (m *Pipeline) Ping() *redis.StatusCmd {
	return m.getPipe().Ping()
}

func (m *Pipeline) Quit() *redis.StatusCmd {
	return m.getPipe().Quit()
}

func (m *Pipeline) Del(keys ...string) *redis.IntCmd {
	return m.getPipe().Del(keys...)
}

func (m *Pipeline) Unlink(keys ...string) *redis.IntCmd {
	return m.getPipe().Unlink(keys...)
}

func (m *Pipeline) Dump(key string) *redis.StringCmd {
	return m.getPipe().Dump(key)
}

func (m *Pipeline) Exists(keys ...string) *redis.IntCmd {
	return m.getPipe().Exists(keys...)
}

func (m *Pipeline) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return m.getPipe().Expire(key, expiration)
}

func (m *Pipeline) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return m.getPipe().PExpireAt(key, tm)
}

func (m *Pipeline) Keys(pattern string) *redis.StringSliceCmd {
	return m.getPipe().Keys(pattern)
}

func (m *Pipeline) Migrate(host, port, key string, db int64, timeout time.Duration) *redis.StatusCmd {
	return m.getPipe().Migrate(host, port, key, db, timeout)
}

func (m *Pipeline) Move(key string, db int64) *redis.BoolCmd {
	return m.getPipe().Move(key, db)
}

func (m *Pipeline) ObjectRefCount(key string) *redis.IntCmd {
	return m.getPipe().ObjectRefCount(key)
}

func (m *Pipeline) ObjectEncoding(key string) *redis.StringCmd {
	return m.getPipe().ObjectEncoding(key)
}

func (m *Pipeline) ObjectIdleTime(key string) *redis.DurationCmd {
	return m.getPipe().ObjectIdleTime(key)
}

func (m *Pipeline) Persist(key string) *redis.BoolCmd {
	return m.getPipe().Persist(key)
}

func (m *Pipeline) PExpire(key string, expiration time.Duration) *redis.BoolCmd {
	return m.getPipe().PExpire(key, expiration)
}

func (m *Pipeline) PExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return m.getPipe().PExpireAt(key, tm)
}
func (m *Pipeline) PTTL(key string) *redis.DurationCmd {
	return m.getPipe().PTTL(key)
}

func (m *Pipeline) RandomKey() *redis.StringCmd {
	return m.getPipe().RandomKey()
}

func (m *Pipeline) Rename(key, newkey string) *redis.StatusCmd {
	return m.getPipe().Rename(key, newkey)
}

func (m *Pipeline) RenameNX(key, newkey string) *redis.BoolCmd {
	return m.getPipe().RenameNX(key, newkey)
}

func (m *Pipeline) Restore(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return m.getPipe().Restore(key, ttl, value)
}

func (m *Pipeline) RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd {
	return m.getPipe().RestoreReplace(key, ttl, value)
}

func (m *Pipeline) Sort(key string, sort *redis.Sort) *redis.StringSliceCmd {
	return m.getPipe().Sort(key, sort)
}

func (m *Pipeline) SortStore(key, store string, sort *redis.Sort) *redis.IntCmd {
	return m.getPipe().SortStore(key, store, sort)
}

func (m *Pipeline) SortInterfaces(key string, sort *redis.Sort) *redis.SliceCmd {
	return m.getPipe().SortInterfaces(key, sort)
}

func (m *Pipeline) Touch(keys ...string) *redis.IntCmd {
	return m.getPipe().Touch(keys...)
}

func (m *Pipeline) TTL(key string) *redis.DurationCmd {
	return m.getPipe().TTL(key)
}

func (m *Pipeline) Type(key string) *redis.StatusCmd {
	return m.getPipe().Type(key)
}

func (m *Pipeline) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return m.getPipe().Scan(cursor, match, count)
}

func (m *Pipeline) SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return m.getPipe().SScan(key, cursor, match, count)
}

func (m *Pipeline) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return m.getPipe().HScan(key, cursor, match, count)
}

func (m *Pipeline) ZScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return m.getPipe().ZScan(key, cursor, match, count)
}

func (m *Pipeline) Append(key, value string) *redis.IntCmd {
	return m.getPipe().Append(key, value)
}

func (m *Pipeline) BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd {
	return m.getPipe().BitCount(key, bitCount)
}

func (m *Pipeline) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	return m.getPipe().BitOpAnd(destKey, keys...)
}

func (m *Pipeline) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	return m.getPipe().BitOpOr(destKey, keys...)
}

func (m *Pipeline) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	return m.getPipe().BitOpXor(destKey, keys...)
}

func (m *Pipeline) BitOpNot(destKey string, key string) *redis.IntCmd {
	return m.getPipe().BitOpNot(destKey, key)
}

func (m *Pipeline) BitPos(key string, bit int64, pos ...int64) *redis.IntCmd {
	return m.getPipe().BitPos(key, bit, pos...)
}

func (m *Pipeline) Decr(key string) *redis.IntCmd {
	return m.getPipe().Decr(key)
}

func (m *Pipeline) DecrBy(key string, decrement int64) *redis.IntCmd {
	return m.getPipe().DecrBy(key, decrement)
}

func (m *Pipeline) Get(key string) *redis.StringCmd {
	return m.getPipe().Get(key)
}

func (m *Pipeline) GetBit(key string, offset int64) *redis.IntCmd {
	return m.getPipe().GetBit(key, offset)
}

func (m *Pipeline) GetRange(key string, start, end int64) *redis.StringCmd {
	return m.getPipe().GetRange(key, start, end)
}

func (m *Pipeline) GetSet(key string, value interface{}) *redis.StringCmd {
	return m.getPipe().GetSet(key, value)
}

func (m *Pipeline) Incr(key string) *redis.IntCmd {
	return m.getPipe().Incr(key)
}

func (m *Pipeline) IncrBy(key string, value int64) *redis.IntCmd {
	return m.getPipe().IncrBy(key, value)
}

func (m *Pipeline) IncrByFloat(key string, value float64) *redis.FloatCmd {
	return m.getPipe().IncrByFloat(key, value)
}

func (m *Pipeline) MGet(keys ...string) *redis.SliceCmd {
	return m.getPipe().MGet(keys...)
}

func (m *Pipeline) MSet(pairs ...interface{}) *redis.StatusCmd {
	return m.getPipe().MSet(pairs...)
}

func (m *Pipeline) MSetNX(pairs ...interface{}) *redis.BoolCmd {
	return m.getPipe().MSetNX(pairs...)
}
func (m *Pipeline) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.getPipe().Set(key, value, expiration)
}

func (m *Pipeline) SetBit(key string, offset int64, value int) *redis.IntCmd {
	return m.getPipe().SetBit(key, offset, value)
}

func (m *Pipeline) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return m.getPipe().SetNX(key, value, expiration)
}

func (m *Pipeline) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return m.getPipe().SetXX(key, value, expiration)
}

func (m *Pipeline) SetRange(key string, offset int64, value string) *redis.IntCmd {
	return m.getPipe().SetRange(key, offset, value)
}

func (m *Pipeline) StrLen(key string) *redis.IntCmd {
	return m.getPipe().StrLen(key)
}

func (m *Pipeline) HDel(key string, fields ...string) *redis.IntCmd {
	return m.getPipe().HDel(key, fields...)
}

func (m *Pipeline) HExists(key, field string) *redis.BoolCmd {
	return m.getPipe().HExists(key, field)
}

func (m *Pipeline) HGet(key, field string) *redis.StringCmd {
	return m.getPipe().HGet(key, field)
}

func (m *Pipeline) HGetAll(key string) *redis.StringStringMapCmd {
	return m.getPipe().HGetAll(key)
}

func (m *Pipeline) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	return m.getPipe().HIncrBy(key, field, incr)
}

func (m *Pipeline) HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	return m.getPipe().HIncrByFloat(key, field, incr)
}

func (m *Pipeline) HKeys(key string) *redis.StringSliceCmd {
	return m.getPipe().HKeys(key)
}

func (m *Pipeline) HLen(key string) *redis.IntCmd {
	return m.getPipe().HLen(key)
}

func (m *Pipeline) HMGet(key string, fields ...string) *redis.SliceCmd {
	return m.getPipe().HMGet(key, fields...)
}

func (m *Pipeline) HMSet(key string, fields map[string]interface{}) *redis.StatusCmd {
	return m.getPipe().HMSet(key, fields)
}

func (m *Pipeline) HSet(key, field string, value interface{}) *redis.BoolCmd {
	return m.getPipe().HSet(key, field, value)
}
func (m *Pipeline) HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	return m.getPipe().HSetNX(key, field, value)
}

func (m *Pipeline) HVals(key string) *redis.StringSliceCmd {
	return m.getPipe().HVals(key)
}

func (m *Pipeline) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return m.getPipe().BLPop(timeout, keys...)
}

func (m *Pipeline) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return m.getPipe().BRPop(timeout, keys...)
}

func (m *Pipeline) BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd {
	return m.getPipe().BRPopLPush(source, destination, timeout)
}

func (m *Pipeline) LIndex(key string, index int64) *redis.StringCmd {
	return m.getPipe().LIndex(key, index)
}

func (m *Pipeline) LInsert(key, op string, pivot, value interface{}) *redis.IntCmd {
	return m.getPipe().LInsert(key, op, pivot, value)
}

func (m *Pipeline) LInsertBefore(key string, pivot, value interface{}) *redis.IntCmd {
	return m.getPipe().LInsertBefore(key, pivot, value)
}

func (m *Pipeline) LInsertAfter(key string, pivot, value interface{}) *redis.IntCmd {
	return m.getPipe().LInsertAfter(key, pivot, value)
}

func (m *Pipeline) LLen(key string) *redis.IntCmd {
	return m.getPipe().LLen(key)
}

func (m *Pipeline) LPop(key string) *redis.StringCmd {
	return m.getPipe().LPop(key)
}

func (m *Pipeline) LPush(key string, values ...interface{}) *redis.IntCmd {
	return m.getPipe().LPush(key, values)
}

func (m *Pipeline) LPushX(key string, value interface{}) *redis.IntCmd {
	return m.getPipe().LPushX(key, value)
}
func (m *Pipeline) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return m.getPipe().LRange(key, start, stop)
}

func (m *Pipeline) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	return m.getPipe().LRem(key, count, value)
}

func (m *Pipeline) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	return m.getPipe().LSet(key, index, value)
}

func (m *Pipeline) LTrim(key string, start, stop int64) *redis.StatusCmd {
	return m.getPipe().LTrim(key, start, stop)
}

func (m *Pipeline) RPop(key string) *redis.StringCmd {
	return m.getPipe().RPop(key)
}

func (m *Pipeline) RPopLPush(source, destination string) *redis.StringCmd {
	return m.getPipe().RPopLPush(source, destination)
}

func (m *Pipeline) RPush(key string, values ...interface{}) *redis.IntCmd {
	return m.getPipe().RPush(key, values...)
}

func (m *Pipeline) RPushX(key string, value interface{}) *redis.IntCmd {
	return m.getPipe().RPushX(key, value)
}

func (m *Pipeline) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return m.getPipe().SAdd(key, members)
}

func (m *Pipeline) SCard(key string) *redis.IntCmd {
	return m.getPipe().SCard(key)
}

func (m *Pipeline) SDiff(keys ...string) *redis.StringSliceCmd {
	return m.getPipe().SDiff(keys...)
}

func (m *Pipeline) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	return m.getPipe().SDiffStore(destination, keys...)
}

func (m *Pipeline) SInter(keys ...string) *redis.StringSliceCmd {
	return m.getPipe().SInter(keys...)
}

func (m *Pipeline) SInterStore(destination string, keys ...string) *redis.IntCmd {
	return m.getPipe().SInterStore(destination, keys...)
}

func (m *Pipeline) SIsMember(key string, member interface{}) *redis.BoolCmd {
	return m.getPipe().SIsMember(key, member)
}

func (m *Pipeline) SMembers(key string) *redis.StringSliceCmd {
	return m.getPipe().SMembers(key)
}

func (m *Pipeline) SMembersMap(key string) *redis.StringStructMapCmd {
	return m.getPipe().SMembersMap(key)
}

func (m *Pipeline) SMove(source, destination string, member interface{}) *redis.BoolCmd {
	return m.getPipe().SMove(source, destination, member)
}

func (m *Pipeline) SPop(key string) *redis.StringCmd {
	return m.getPipe().SPop(key)
}

func (m *Pipeline) SPopN(key string, count int64) *redis.StringSliceCmd {
	return m.getPipe().SPopN(key, count)
}

func (m *Pipeline) SRandMember(key string) *redis.StringCmd {
	return m.getPipe().SRandMember(key)
}

func (m *Pipeline) SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	return m.getPipe().SRandMemberN(key, count)
}

func (m *Pipeline) SRem(key string, members ...interface{}) *redis.IntCmd {
	return m.getPipe().SRem(key, members)
}

func (m *Pipeline) SUnion(keys ...string) *redis.StringSliceCmd {
	return m.getPipe().SUnion(keys...)
}

func (m *Pipeline) SUnionStore(destination string, keys ...string) *redis.IntCmd {
	return m.getPipe().SUnionStore(destination, keys...)
}

func (m *Pipeline) XAdd(a *redis.XAddArgs) *redis.StringCmd {
	return m.getPipe().XAdd(a)
}

func (m *Pipeline) XDel(stream string, ids ...string) *redis.IntCmd {
	return m.getPipe().XDel(stream, ids...)
}

func (m *Pipeline) XLen(stream string) *redis.IntCmd {
	return m.getPipe().XLen(stream)
}

func (m *Pipeline) XRange(stream, start, stop string) *redis.XMessageSliceCmd {
	return m.getPipe().XRange(stream, start, stop)
}

func (m *Pipeline) XRangeN(stream, start, stop string, count int64) *redis.XMessageSliceCmd {
	return m.getPipe().XRangeN(stream, start, stop, count)
}

func (m *Pipeline) XRevRange(stream string, start, stop string) *redis.XMessageSliceCmd {
	return m.getPipe().XRevRange(stream, start, stop)
}

func (m *Pipeline) XRevRangeN(stream string, start, stop string, count int64) *redis.XMessageSliceCmd {
	return m.getPipe().XRevRangeN(stream, start, stop, count)
}

func (m *Pipeline) XRead(a *redis.XReadArgs) *redis.XStreamSliceCmd {
	return m.getPipe().XRead(a)
}

func (m *Pipeline) XReadStreams(streams ...string) *redis.XStreamSliceCmd {
	return m.getPipe().XReadStreams(streams...)
}

func (m *Pipeline) XGroupCreate(stream, group, start string) *redis.StatusCmd {
	return m.getPipe().XGroupCreate(stream, group, start)
}

func (m *Pipeline) XGroupCreateMkStream(stream, group, start string) *redis.StatusCmd {
	return m.getPipe().XGroupCreateMkStream(stream, group, start)
}

func (m *Pipeline) XGroupSetID(stream, group, start string) *redis.StatusCmd {
	return m.getPipe().XGroupSetID(stream, group, start)
}

func (m *Pipeline) XGroupDestroy(stream, group string) *redis.IntCmd {
	return m.getPipe().XGroupDestroy(stream, group)
}

func (m *Pipeline) XGroupDelConsumer(stream, group, consumer string) *redis.IntCmd {
	return m.getPipe().XGroupDelConsumer(stream, group, consumer)
}

func (m *Pipeline) XReadGroup(a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	return m.getPipe().XReadGroup(a)
}

func (m *Pipeline) XAck(stream, group string, ids ...string) *redis.IntCmd {
	return m.getPipe().XAck(stream, group, ids...)
}

func (m *Pipeline) XPending(stream, group string) *redis.XPendingCmd {
	return m.getPipe().XPending(stream, group)
}

func (m *Pipeline) XPendingExt(a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	return m.getPipe().XPendingExt(a)
}

func (m *Pipeline) XClaim(a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	return m.getPipe().XClaim(a)
}

func (m *Pipeline) XClaimJustID(a *redis.XClaimArgs) *redis.StringSliceCmd {
	return m.getPipe().XClaimJustID(a)
}

func (m *Pipeline) XTrim(key string, maxLen int64) *redis.IntCmd {
	return m.getPipe().XTrim(key, maxLen)
}

func (m *Pipeline) XTrimApprox(key string, maxLen int64) *redis.IntCmd {
	return m.getPipe().XTrimApprox(key, maxLen)
}

func (m *Pipeline) BZPopMax(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return m.getPipe().BZPopMax(timeout, keys...)
}

func (m *Pipeline) BZPopMin(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	return m.getPipe().BZPopMin(timeout, keys...)
}
func (m *Pipeline) ZAdd(key string, members ...redis.Z) *redis.IntCmd {
	return m.getPipe().ZAdd(key, members...)
}

func (m *Pipeline) ZAddNX(key string, members ...redis.Z) *redis.IntCmd {
	return m.getPipe().ZAddNX(key, members...)
}

func (m *Pipeline) ZAddXX(key string, members ...redis.Z) *redis.IntCmd {
	return m.getPipe().ZAddXX(key, members...)
}

func (m *Pipeline) ZAddCh(key string, members ...redis.Z) *redis.IntCmd {
	return m.getPipe().ZAddCh(key, members...)
}

func (m *Pipeline) ZAddNXCh(key string, members ...redis.Z) *redis.IntCmd {
	return m.getPipe().ZAddNXCh(key, members...)
}

func (m *Pipeline) ZAddXXCh(key string, members ...redis.Z) *redis.IntCmd {
	return m.getPipe().ZAddXXCh(key, members...)
}

func (m *Pipeline) ZIncr(key string, member redis.Z) *redis.FloatCmd {
	return m.getPipe().ZIncr(key, member)
}

func (m *Pipeline) ZIncrNX(key string, member redis.Z) *redis.FloatCmd {
	return m.getPipe().ZIncrNX(key, member)
}

func (m *Pipeline) ZIncrXX(key string, member redis.Z) *redis.FloatCmd {
	return m.getPipe().ZIncrXX(key, member)
}

func (m *Pipeline) ZCard(key string) *redis.IntCmd {
	return m.getPipe().ZCard(key)
}

func (m *Pipeline) ZCount(key, min, max string) *redis.IntCmd {
	return m.getPipe().ZCount(key, min, max)
}

func (m *Pipeline) ZLexCount(key, min, max string) *redis.IntCmd {
	return m.getPipe().ZLexCount(key, min, max)
}

func (m *Pipeline) ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	return m.getPipe().ZIncrBy(key, increment, member)
}

func (m *Pipeline) ZInterStore(destination string, store redis.ZStore, keys ...string) *redis.IntCmd {
	return m.getPipe().ZInterStore(destination, store, keys...)
}

func (m *Pipeline) ZPopMax(key string, count ...int64) *redis.ZSliceCmd {
	return m.getPipe().ZPopMax(key, count...)
}

func (m *Pipeline) ZPopMin(key string, count ...int64) *redis.ZSliceCmd {
	return m.getPipe().ZPopMin(key, count...)
}

func (m *Pipeline) ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	return m.getPipe().ZRange(key, start, stop)
}

func (m *Pipeline) ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return m.getPipe().ZRangeWithScores(key, start, stop)
}

func (m *Pipeline) ZRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return m.getPipe().ZRangeByScore(key, opt)
}

func (m *Pipeline) ZRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return m.getPipe().ZRangeByLex(key, opt)
}

func (m *Pipeline) ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	return m.getPipe().ZRangeByScoreWithScores(key, opt)
}

func (m *Pipeline) ZRank(key, member string) *redis.IntCmd {
	return m.getPipe().ZRank(key, member)
}

func (m *Pipeline) ZRem(key string, members ...interface{}) *redis.IntCmd {
	return m.getPipe().ZRem(key, members...)
}

func (m *Pipeline) ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd {
	return m.getPipe().ZRemRangeByRank(key, start, stop)
}

func (m *Pipeline) ZRemRangeByScore(key, min, max string) *redis.IntCmd {
	return m.getPipe().ZRemRangeByScore(key, min, max)
}

func (m *Pipeline) ZRemRangeByLex(key, min, max string) *redis.IntCmd {
	return m.getPipe().ZRemRangeByLex(key, min, max)
}

func (m *Pipeline) ZRevRange(key string, start, stop int64) *redis.StringSliceCmd {
	return m.getPipe().ZRevRange(key, start, stop)
}

func (m *Pipeline) ZRevRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return m.getPipe().ZRevRangeWithScores(key, start, stop)
}

func (m *Pipeline) ZRevRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return m.getPipe().ZRevRangeByScore(key, opt)
}

func (m *Pipeline) ZRevRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	return m.getPipe().ZRangeByLex(key, opt)
}

func (m *Pipeline) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	return m.getPipe().ZRangeByScoreWithScores(key, opt)
}

func (m *Pipeline) ZRevRank(key, member string) *redis.IntCmd {
	return m.getPipe().ZRevRank(key, member)
}

func (m *Pipeline) ZScore(key, member string) *redis.FloatCmd {
	return m.getPipe().ZScore(key, member)
}

func (m *Pipeline) ZUnionStore(dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	return m.getPipe().ZUnionStore(dest, store, keys...)
}

func (m *Pipeline) PFAdd(key string, els ...interface{}) *redis.IntCmd {
	return m.getPipe().PFAdd(key, els...)
}

func (m *Pipeline) PFCount(keys ...string) *redis.IntCmd {
	return m.getPipe().PFCount(keys...)
}

func (m *Pipeline) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	return m.getPipe().PFMerge(dest, keys...)
}

func (m *Pipeline) BgRewriteAOF() *redis.StatusCmd {
	return m.getPipe().BgRewriteAOF()
}

func (m *Pipeline) BgSave() *redis.StatusCmd {
	return m.getPipe().BgSave()
}

func (m *Pipeline) ClientKill(ipPort string) *redis.StatusCmd {
	return m.getPipe().ClientKill(ipPort)
}

func (m *Pipeline) ClientKillByFilter(keys ...string) *redis.IntCmd {
	return m.getPipe().ClientKillByFilter(keys...)
}

func (m *Pipeline) ClientList() *redis.StringCmd {
	return m.getPipe().ClientList()
}

func (m *Pipeline) ClientPause(dur time.Duration) *redis.BoolCmd {
	return m.getPipe().ClientPause(dur)
}

func (m *Pipeline) ClientID() *redis.IntCmd {
	return m.getPipe().ClientID()
}

func (m *Pipeline) ConfigGet(parameter string) *redis.SliceCmd {
	return m.getPipe().ConfigGet(parameter)
}

func (m *Pipeline) ConfigResetStat() *redis.StatusCmd {
	return m.getPipe().ConfigResetStat()
}

func (m *Pipeline) ConfigSet(parameter, value string) *redis.StatusCmd {
	return m.getPipe().ConfigSet(parameter, value)
}

func (m *Pipeline) ConfigRewrite() *redis.StatusCmd {
	return m.getPipe().ConfigRewrite()
}

func (m *Pipeline) DBSize() *redis.IntCmd {
	return m.getPipe().DBSize()
}

func (m *Pipeline) FlushAll() *redis.StatusCmd {
	return m.getPipe().FlushAll()
}
func (m *Pipeline) FlushAllAsync() *redis.StatusCmd {
	return m.getPipe().FlushAllAsync()
}

func (m *Pipeline) FlushDB() *redis.StatusCmd {
	return m.getPipe().FlushDB()
}

func (m *Pipeline) FlushDBAsync() *redis.StatusCmd {
	return m.getPipe().FlushDBAsync()
}

func (m *Pipeline) Info(section ...string) *redis.StringCmd {
	return m.getPipe().Info(section...)
}

func (m *Pipeline) LastSave() *redis.IntCmd {
	return m.getPipe().LastSave()
}

func (m *Pipeline) Save() *redis.StatusCmd {
	return m.getPipe().Save()
}

func (m *Pipeline) Shutdown() *redis.StatusCmd {
	return m.getPipe().Shutdown()
}

func (m *Pipeline) ShutdownSave() *redis.StatusCmd {
	return m.getPipe().ShutdownSave()
}

func (m *Pipeline) ShutdownNoSave() *redis.StatusCmd {
	return m.getPipe().ShutdownNoSave()
}

func (m *Pipeline) SlaveOf(host, port string) *redis.StatusCmd {
	return m.getPipe().SlaveOf(host, port)
}

func (m *Pipeline) Time() *redis.TimeCmd {
	return m.getPipe().Time()
}

func (m *Pipeline) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	return m.getPipe().Eval(script, keys, args...)
}

func (m *Pipeline) EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	return m.getPipe().EvalSha(sha1, keys, args...)
}

func (m *Pipeline) ScriptExists(hashes ...string) *redis.BoolSliceCmd {
	return m.getPipe().ScriptExists(hashes...)
}

func (m *Pipeline) ScriptFlush() *redis.StatusCmd {
	return m.getPipe().ScriptFlush()
}

func (m *Pipeline) ScriptKill() *redis.StatusCmd {
	return m.getPipe().ScriptKill()
}

func (m *Pipeline) ScriptLoad(script string) *redis.StringCmd {
	return m.getPipe().ScriptLoad(script)
}

func (m *Pipeline) DebugObject(key string) *redis.StringCmd {
	return m.getPipe().DebugObject(key)
}

func (m *Pipeline) Publish(channel string, message interface{}) *redis.IntCmd {
	return m.getPipe().Publish(channel, message)
}

func (m *Pipeline) PubSubChannels(pattern string) *redis.StringSliceCmd {
	return m.getPipe().PubSubChannels(pattern)
}

func (m *Pipeline) PubSubNumSub(channels ...string) *redis.StringIntMapCmd {
	return m.getPipe().PubSubNumSub(channels...)
}

func (m *Pipeline) PubSubNumPat() *redis.IntCmd {
	return m.getPipe().PubSubNumPat()
}

func (m *Pipeline) ClusterSlots() *redis.ClusterSlotsCmd {
	return m.getPipe().ClusterSlots()
}

func (m *Pipeline) ClusterNodes() *redis.StringCmd {
	return m.getPipe().ClusterNodes()
}

func (m *Pipeline) ClusterMeet(host, port string) *redis.StatusCmd {
	return m.getPipe().ClusterMeet(host, port)
}

func (m *Pipeline) ClusterForget(nodeID string) *redis.StatusCmd {
	return m.getPipe().ClusterForget(nodeID)
}

func (m *Pipeline) ClusterReplicate(nodeID string) *redis.StatusCmd {
	return m.getPipe().ClusterReplicate(nodeID)
}

func (m *Pipeline) ClusterResetSoft() *redis.StatusCmd {
	return m.getPipe().ClusterResetSoft()
}

func (m *Pipeline) ClusterResetHard() *redis.StatusCmd {
	return m.getPipe().ClusterResetHard()
}

func (m *Pipeline) ClusterInfo() *redis.StringCmd {
	return m.getPipe().ClusterInfo()
}

func (m *Pipeline) ClusterKeySlot(key string) *redis.IntCmd {
	return m.getPipe().ClusterKeySlot(key)
}

func (m *Pipeline) ClusterGetKeysInSlot(slot int, count int) *redis.StringSliceCmd {
	return m.getPipe().ClusterGetKeysInSlot(slot, count)
}

func (m *Pipeline) ClusterCountFailureReports(nodeID string) *redis.IntCmd {
	return m.getPipe().ClusterCountFailureReports(nodeID)
}

func (m *Pipeline) ClusterCountKeysInSlot(slot int) *redis.IntCmd {
	return m.getPipe().ClusterCountKeysInSlot(slot)
}

func (m *Pipeline) ClusterDelSlots(slots ...int) *redis.StatusCmd {
	return m.getPipe().ClusterDelSlots(slots...)
}

func (m *Pipeline) ClusterDelSlotsRange(min, max int) *redis.StatusCmd {
	return m.getPipe().ClusterDelSlotsRange(min, max)
}

func (m *Pipeline) ClusterSaveConfig() *redis.StatusCmd {
	return m.getPipe().ClusterSaveConfig()
}

func (m *Pipeline) ClusterSlaves(nodeID string) *redis.StringSliceCmd {
	return m.getPipe().ClusterSlaves(nodeID)
}

func (m *Pipeline) ClusterFailover() *redis.StatusCmd {
	return m.getPipe().ClusterFailover()
}

func (m *Pipeline) ClusterAddSlots(slots ...int) *redis.StatusCmd {
	return m.getPipe().ClusterAddSlots(slots...)
}

func (m *Pipeline) ClusterAddSlotsRange(min, max int) *redis.StatusCmd {
	return m.getPipe().ClusterAddSlotsRange(min, max)
}

func (m *Pipeline) GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	return m.getPipe().GeoAdd(key, geoLocation...)
}

func (m *Pipeline) GeoPos(key string, members ...string) *redis.GeoPosCmd {
	return m.getPipe().GeoPos(key, members...)
}

func (m *Pipeline) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return m.getPipe().GeoRadius(key, longitude, latitude, query)
}

func (m *Pipeline) GeoRadiusRO(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return m.getPipe().GeoRadiusRO(key, longitude, latitude, query)
}

func (m *Pipeline) GeoRadiusByMember(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return m.getPipe().GeoRadiusByMember(key, member, query)
}

func (m *Pipeline) GeoRadiusByMemberRO(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	return m.getPipe().GeoRadiusByMemberRO(key, member, query)
}

func (m *Pipeline) GeoDist(key string, member1, member2, unit string) *redis.FloatCmd {
	return m.getPipe().GeoDist(key, member1, member2, unit)
}

func (m *Pipeline) GeoHash(key string, members ...string) *redis.StringSliceCmd {
	return m.getPipe().GeoHash(key, members...)
}

func (m *Pipeline) ReadOnly() *redis.StatusCmd {
	return m.getPipe().ReadOnly()
}

func (m *Pipeline) ReadWrite() *redis.StatusCmd {
	return m.getPipe().ReadWrite()
}

func (m *Pipeline) MemoryUsage(key string, samples ...int) *redis.IntCmd {
	return m.getPipe().MemoryUsage(key, samples...)
}
