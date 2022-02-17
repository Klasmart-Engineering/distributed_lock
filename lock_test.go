package lock

import (
	"context"
	"fmt"
	"gitlab.badanamu.com.cn/calmisland/distributed_lock/drivers"
	"gitlab.badanamu.com.cn/calmisland/distributed_lock/utils"
	"strconv"
	"sync"
	"testing"
	"time"
)
var wg sync.WaitGroup

func testPrint(ld LockDriver, s string){
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
			Host:    "192.168.1.234",
			Port:     6379,
			Password: "",
		},
	})
	if err != nil{
		panic(err)
	}

	for i := 0; i < 100; i ++ {
		wg.Add(1)
		go testPrint(redisLock, strconv.Itoa(i))
		//time.Sleep(time.Second)
	}

	wg.Wait()
}

func TestLock0(t *testing.T) {
	x := 0
	cfg := DistributedLockConfig{
		Driver:  "redis",
		Key:     fmt.Sprintf("lock.order.%v", utils.RandNum()),
		Timeout: time.Minute * 1,
		Ctx:     context.Background(),
		RedisConfig: drivers.RedisConfig{
			Host:    "192.168.1.234",
			Port:     6379,
			Password: "",
		},
	}

	wg := new(sync.WaitGroup)
	for i := 0; i < 1000; i ++ {
		wg.Add(1)
		go func() {
			redisLock, err := NewRedisLock(cfg)
			if err != nil{
				panic(err)
			}
			redisLock.Lock()
			x ++
			redisLock.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(x)
}