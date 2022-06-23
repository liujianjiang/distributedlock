package distributedlock

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

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

func TestRedisSpinLock(t *testing.T) {
	l := NewRedisLock(rediscli, "test-key-TestRedisSpinLock")
	success, err := l.SpinLock(context.Background())
	assert.Nil(t, err)
	assert.True(t, success)
	err = l.Unlock(context.Background())
	assert.NotNil(t, err)
}

func TestRedisWatchLock(t *testing.T) {
	l := NewRedisLock(rediscli, "test-key-TestRedisWatchLock", WithTTL(time.Millisecond*10000))
	success, err := l.Lock(context.Background())
	assert.Nil(t, err)
	assert.True(t, success)

	ctx, cancelFun := context.WithCancel(context.Background())
	//开启协程监听锁并自动续约
	go l.WatchLock(ctx)
	//处理业务
	if success {
		time.Sleep(time.Second * 40)
	}
	//业务处理完成关闭监听协程
	cancelFun()
	err = l.Unlock(context.Background())
	assert.Nil(t, err)
}
