package lock

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/KL-Engineering/distributed_lock/drivers"
)

var wg sync.WaitGroup

func testPrint(ld LockDriver, s string) {
	ld.Lock()
	defer func() {
		ld.Unlock()
		fmt.Println("Done")
	}()
	fmt.Println("Greet - " + s)
	wg.Done()
}

func TestLock(t *testing.T) {
	redisLock, err := NewRedisLock(DistributedLockConfig{
		Driver:  "redis",
		Key:     fmt.Sprintf("lock.order.112"),
		Timeout: time.Minute * 1,
		Ctx:     context.Background(),
		RedisConfig: drivers.RedisConfig{
			Host:     "127.0.0.1",
			Port:     6379,
			Password: "",
		},
	})
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go testPrint(redisLock, strconv.Itoa(i))
		//time.Sleep(time.Second)
	}

	wg.Wait()
}
