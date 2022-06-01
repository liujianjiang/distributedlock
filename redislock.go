// Package distributedlock 提供分布式锁的统一实现
package distributedlock

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisLock struct {
	cli   redis.UniversalClient
	key   string
	value interface{}
	ttl   time.Duration
}

var _ DistributedLock = &redisLock{}

// lua 脚本加锁和过期时间
var redisLockLua = redis.NewScript(`
local key = KEYS[1]
local value = ARGV[1]
local ttl = tonumber(ARGV[2])
local ret = redis.call("setnx", key, value)
if ret == 1 then
	redis.call("pexpire", key, ttl)
else
	if redis.call("get", key) == value then
		ret = 1
		redis.call("pexpire", key, ttl)
	end
end
return ret`)

// lua脚本释放锁
var redisUnlockLua = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`)

// NewRedisLock 基于 Redis 实现，使用 go-redis
func NewRedisLock(cli redis.UniversalClient, key string, opts ...Option) DistributedLock {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	l := &redisLock{
		cli:   cli,
		key:   key,
		value: options.Value,
		ttl:   options.TTL,
	}
	return l
}

// 加锁
func (l *redisLock) Lock(ctx context.Context) (bool, error) {
	ttl := strconv.FormatInt(int64(l.ttl/time.Millisecond), 10)
	status, err := redisLockLua.Run(ctx, l.cli, []string{l.key}, l.value, ttl).Result()
	if err != nil {
		return false, err
	}
	if i, ok := status.(int64); ok && i == 1 {
		return true, nil
	}
	return false, nil
}

// 自旋锁
func (l *redisLock) SpinLock(ctx context.Context) (success bool, err error) {
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}
		success, err = l.Lock(ctx)
		fmt.Printf("err: %v", err)
		if err != nil {

			time.Sleep(10 * time.Millisecond)
			continue
		}
		if success {
			return
		}
	}
}

// 释放锁
func (l *redisLock) Unlock(ctx context.Context) error {
	status, err := redisUnlockLua.Run(ctx, l.cli, []string{l.key}, l.value).Result()
	if err == redis.Nil {
		return ErrLockNotHeld
	} else if err != nil {
		return err
	}
	if i, ok := status.(int64); !ok || i != 1 {
		return ErrLockNotHeld
	}
	return nil
}
