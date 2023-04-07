package redislock_test

import (
	"context"
	"time"

	"github.com/khasanovbi/redislock"
	"github.com/redis/go-redis/v9"
)

func doJobUnderLockWithContext(ctx context.Context) { //nolint:revive
}

func ExampleLocker_Lock() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	locker, err := redislock.New(rdb, redislock.WithRefreshPeriod(time.Minute), redislock.WithTTL(10*time.Minute))
	if err != nil {
		panic(err)
	}

	lock, ctx, err := locker.Lock(context.Background(), "my_lock_key")
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = lock.Unlock()
	}()

	doJobUnderLockWithContext(ctx)
}
