// Package distributedlock 提供分布式锁的统一实现
// 要求：所有实现的 Lock/SpinLock 需要支持可重入
package distributedlock

import (
	"context"
	"errors"
)

// DistributedLock 分布式锁
type DistributedLock interface {
	// Lock 获取锁，非阻塞，能不能获取到都立即返回获取结果
	Lock(ctx context.Context) (success bool, err error)

	// SpinLock 获取锁，自旋，获取不到就阻塞等待，直到 context 超时
	SpinLock(ctx context.Context) (success bool, err error)

	// WatchLock 创建守护协程，监听锁自动续期
	WatchLock(ctx context.Context) error

	// Unlock 释放锁
	Unlock(ctx context.Context) error
}

var (
	ErrLockNotHeld = errors.New("distributedlock: lock not held")
)
