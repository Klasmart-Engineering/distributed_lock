package lock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KL-Engineering/distributed_lock/drivers"
)

var (
	ErrUnknownDistributedLockDriver = errors.New("Unknown distributed lock driver")
	ErrLockTimeout                  = errors.New("Lock timeout")
)

type LockDriver interface {
	Lock()
	Unlock()
}

type RedisLock struct {
	dc          DistributedLockConfig
	lockChannel chan bool
	exitChannel chan struct{}
	isGetLock   bool
}

func (r *RedisLock) Lock() {
	//尝试等待5s
	go func() {
		for {
			select {
			case <-r.exitChannel:
				break
			default:
				ret, _ := drivers.GetRedis().SetNX(r.dc.Key, "1", r.dc.Timeout).Result()
				if ret {
					fmt.Println("Already get lock")
					r.isGetLock = true
					r.lockChannel <- true
					return
				}
				time.Sleep(time.Duration(10) * time.Millisecond)
			}

		}
	}()

	select {
	case <-r.lockChannel:
		fmt.Println("Got lock")
		//r.exitChannel <- struct{}{}
		return
	case <-r.dc.Ctx.Done():
		r.exitChannel <- struct{}{}
		return
	}

}

func (r *RedisLock) Unlock() {
	fmt.Println("Unlock Res:", r.isGetLock)
	if !r.isGetLock {
		return
	}
	drivers.GetRedis().Del(r.dc.Key)
	r.isGetLock = false
}

func NewRedisLock(dc DistributedLockConfig) (LockDriver, error) {
	err := drivers.OpenRedis(dc.RedisConfig)
	if err != nil {
		return nil, err
	}
	return &RedisLock{
		dc:          dc,
		lockChannel: make(chan bool),
		exitChannel: make(chan struct{}),
	}, nil
}

type DistributedLockConfig struct {
	Driver      string
	Key         string
	Timeout     time.Duration
	Ctx         context.Context
	RedisConfig drivers.RedisConfig
}

func NewDistributedLock(dc DistributedLockConfig) (LockDriver, error) {
	if dc.Driver == "redis" {
		//打开redis
		return NewRedisLock(dc)
	}
	return nil, ErrUnknownDistributedLockDriver
}
