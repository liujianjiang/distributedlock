# distributedlock 分布式锁

```go
import "github.com/liujianjiang/distributedlock"

// Hello is an example method
func main() {
	lockkey := "my-business-lock-key"
	l := distributedlock.NewRedisLock(s.rediscli, lockkey)
	// Lock 已经支持可重入，不用担心 Redis中设置成功而客户端超时的场景下重试后获取不到锁的情况
	success, err := l.Lock(ctx)
	if err != nil {
		// err != nil 的时候说明有系统错误，业务可选择重试获取 或 放弃获取
		log.ErrorContextf("Hello: lock key(%s) error(%v)", lockkey, err)
		return err
	}
	if !success {
		// 处理没获取到锁的情况
		return nil
	}
	// 这里已经获取到了锁，可以处理其他业务逻辑，处理完后释放锁
	if err = l.Unlock(ctx); err != nil {
		log.ErrorContextf("Hello: unlock key(%s) error(%v)", lockkey, err)
		return err
	}
	return nil
}
```
