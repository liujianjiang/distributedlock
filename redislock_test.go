package distributedlock

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var rediscli redis.UniversalClient

func TestMain(m *testing.M) {
	rediscli = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	os.Exit(m.Run())
}

// TestRedisLock 测试基础的 Lock & Unlock
func TestRedisLock(t *testing.T) {
	l := NewRedisLock(rediscli, "test-key-TestRedisLock")
	success, err := l.Lock(context.Background())
	fmt.Println(err)
	assert.Nil(t, err)
	assert.True(t, success)
	err = l.Unlock(context.Background())
	assert.Nil(t, err)
	err = l.Unlock(context.Background())
	assert.NotNil(t, err)
}

// TestRedisReentrantLock 测试可重入锁
func TestRedisReentrantLock(t *testing.T) {
	l := NewRedisLock(rediscli, "test-key-TestRedisReentrantLock")
	success, err := l.Lock(context.Background())
	assert.Nil(t, err)
	assert.True(t, success)
	success, err = l.Lock(context.Background())
	assert.Nil(t, err)
	assert.True(t, success)
	err = l.Unlock(context.Background())
	assert.Nil(t, err)
	err = l.Unlock(context.Background())
	assert.NotNil(t, err)
}
