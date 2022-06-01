package distributedlock

import (
	"math/rand"
	"time"
)

const randCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// DefaultValueLength 默认锁值的长度
var DefaultValueLength = 8

// Options is lock options
type Options struct {
	// Value 是锁的值，不同的实现可能对 Value 使用方式不同，也可能不会使用 Value
	// 在 redislock 实现中，Value 存放每把锁的特定标识，用来防止客户端A释放了客户端B的锁
	Value interface{}

	// TTL 是锁的时间
	TTL time.Duration
}

// Option is lock option func
type Option func(*Options)

// DefaultOptions 默认锁配置
func DefaultOptions() *Options {
	return &Options{
		Value: randstring(DefaultValueLength),
		TTL:   5 * time.Second,
	}
}

func randstring(length int) string {
	s := make([]byte, length)
	l := len(randCharset)
	for i := range s {
		s[i] = randCharset[rand.Intn(l)]
	}
	return string(s)
}

// WithValue sets lock value
func WithValue(val interface{}) Option {
	return func(opts *Options) {
		opts.Value = val
	}
}

// WithTTL sets lock duration
func WithTTL(ttl time.Duration) Option {
	return func(opts *Options) {
		opts.TTL = ttl
	}
}
