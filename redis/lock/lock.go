package lock

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-redis/redis/v8"
)

const (
	// 过期时间
	defaultTTL = time.Second * 3
	// 重置过期时间间隔
	defaultResetInterval = defaultTTL * 2 / 3
	// 重新获取锁间隔
	defaultTryLockInterval = time.Microsecond * 500
	// 默认重试次数
	defaultMaxAttempt = 3
	// 解锁脚本
	unlock_script = `
		if redis.call("get",KEYS[1]) == ARGV[1] then
			return redis.call("del",KEYS[1])
		else
			return 0
		end
	`
	// 续期脚本
	expire_script = `
		if redis.call("get",KEYS[1]) == ARGV[1] then
			return redis.call("expire",KEYS[1],ARGV[2])
		else
			return 0
		end
	`
)

var (
	ErrLockFailed    = errors.New("lock failed")
	ErrTimeout       = errors.New("timeout")
	ErrMaxAttempt    = errors.New("exceeded maximum attempts")
	ErrKeyNotMatched = errors.New("key value not matched")
)

type Option func(*Lock)

func WithTTL(ttl time.Duration) Option {
	return func(lock *Lock) {
		lock.ttl = ttl
	}
}

func WithTryLockInterval(tryLockInterval time.Duration) Option {
	return func(lock *Lock) {
		lock.tryLockInterval = tryLockInterval
	}
}

func WithResetInterval(resetInterval time.Duration) Option {
	return func(lock *Lock) {
		lock.resetInterval = resetInterval
	}
}

func WithRenew(renew bool) Option {
	return func(lock *Lock) {
		lock.renew = renew
	}
}

func WithMaxAttemp(maxAttemp int) Option {
	return func(lock *Lock) {
		lock.maxAttemp = maxAttemp
	}
}

// Lock 不可重复使用
type Lock struct {
	client       *redis.Client // Redis客户端
	key          string        // key
	value        string        // value
	unlockScript *redis.Script // 解锁脚本
	expireScript *redis.Script // 续期脚本
	watchDog     chan struct{} // 看门狗

	ttl             time.Duration // 过期时间
	tryLockInterval time.Duration // 重新获取锁间隔
	resetInterval   time.Duration // 重制续期时间
	maxAttemp       int           // 最大尝试次数
	renew           bool          // 是否开启自动续期
}

func NewLock(client *redis.Client, key string, options ...Option) *Lock {
	lock := &Lock{
		client:       client,
		unlockScript: redis.NewScript(unlock_script),
		expireScript: redis.NewScript(expire_script),
		watchDog:     make(chan struct{}),
		key:          key,
		value:        gofakeit.UUID(),
	}

	for _, opt := range options {
		opt(lock)
	}

	if lock.maxAttemp <= 0 {
		lock.maxAttemp = defaultMaxAttempt
	}

	if lock.ttl <= 0 {
		lock.ttl = defaultTTL
	}

	if lock.tryLockInterval <= 0 {
		lock.tryLockInterval = defaultTryLockInterval
	}

	if lock.resetInterval <= 0 {
		lock.resetInterval = defaultResetInterval
	}

	return lock
}

func (l *Lock) Lock(ctx context.Context) error {
	// 尝试加锁
	if err := l.TryLock(ctx); err == nil {
		return nil
	}
	// 加锁失败，不断尝试
	ticker := time.NewTicker(l.tryLockInterval)
	defer ticker.Stop()
	for i := 0; i < l.maxAttemp; i++ {
		select {
		case <-ctx.Done():
			// 超时
			return ErrTimeout
		case <-ticker.C:
			// 重新尝试加锁
			if err := l.TryLock(ctx); err == nil {
				return nil // 加锁成功
			}
		}
	}
	return ErrMaxAttempt
}

func (l *Lock) TryLock(ctx context.Context) error {
	success, err := l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
	if err != nil {
		return err
	}
	// 加锁失败
	if !success {
		return ErrLockFailed
	}
	// 加锁成功，启动看门狗
	go l.startWatchDog(ctx)
	return nil
}

func (l *Lock) startWatchDog(ctx context.Context) {
	ticker := time.NewTicker(l.resetInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			result, err := l.expireScript.Run(ctx, l.client, []string{l.key}, l.value, int(l.ttl.Seconds())).Result()
			if err != nil {
				log.Printf("watch dog fail for renew key: %s\n", l.key)
				return
			}

			if result.(int64) == 0 {
				log.Printf("watch dog fail for key value not matched: %s\n", l.key)
				return
			}

			log.Printf("watch dog succ for key: %s\n", l.key)
		case <-l.watchDog:
			// 已经解锁
			log.Printf("watch dog exit for key: %s\n", l.key)
			return
		}
	}
}

func (l *Lock) Unlock(ctx context.Context) error {
	result, err := l.unlockScript.Run(ctx, l.client, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}
	if result.(int64) == 0 {
		return ErrKeyNotMatched
	}
	close(l.watchDog)
	return err
}

type LockManager struct {
	client *redis.Client // Redis客户端
}

func NewLockManager(client *redis.Client) *LockManager {
	return &LockManager{
		client: client,
	}
}

func (m *LockManager) CreateLock(key string, options ...Option) *Lock {
	lock := NewLock(m.client, key, options...)
	return lock
}
