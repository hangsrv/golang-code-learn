package lock

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestNewLock(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	lm := NewLockManager(redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
	}))

	var count = 5
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		locker := lm.CreateLock("test", WithTryLockInterval(100*time.Millisecond))
		go func(locker *Lock, num int) {
			defer wg.Done()
			err := locker.Lock(context.Background())
			if err != nil {
				if err == ErrMaxAttempt {
					log.Println(err.Error())
					return
				}
				log.Println(err.Error())
				t.Error(err)
			}
			for i := 0; i < 10; i++ {
				log.Printf("business:%d\n", num)
				time.Sleep(time.Second * 1)
			}

			err = locker.Unlock(context.Background())
			if err != nil {
				log.Println(err.Error())
				t.Error(err)
			}
		}(locker, i)
	}
	wg.Wait()
}
