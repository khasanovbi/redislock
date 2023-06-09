# redislock
Redis context-based locker with auto refresh.

## Examples

```go
package main

import (
	"context"
	"time"

	"github.com/khasanovbi/redislock"
	"github.com/redis/go-redis/v9"
)

func doJobUnderLockWithContext(ctx context.Context) {
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Lock will automatically update the TTL time every RefreshPeriod.
	// The configuration can be set either for the entire Locker or for each Lock request individually.
	locker, err := redislock.New(rdb, redislock.WithRefreshPeriod(time.Minute), redislock.WithTTL(10*time.Minute))
	if err != nil {
		panic(err)
	}

	// Try to get lock.
	lock, ctx, err := locker.Lock(context.Background(), "my_lock_key")
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = lock.Unlock()
	}()
	
	doJobUnderLockWithContext(ctx)
}
```