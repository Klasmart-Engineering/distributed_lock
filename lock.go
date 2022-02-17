package lock

import (
	"context"
	"errors"
	"gitlab.badanamu.com.cn/calmisland/distributed_lock/drivers"
	"gitlab.badanamu.com.cn/calmisland/distributed_lock/utils"
	"time"
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
	value	string
}

func (r *RedisLock) Lock() {
	go func() {
		for {
			select {
			case <-r.exitChannel:
				return
			default:
				ret, _ := drivers.GetRedis().SetNX(r.dc.Key, r.value, r.dc.Timeout).Result()
				if ret {
					r.lockChannel <- true
					return
				}
				time.Sleep(time.Duration(60) * time.Millisecond)
			}

		}
	}()

	select {
	case <-r.lockChannel:
		return
	case <-r.dc.Ctx.Done():
		r.exitChannel <- struct{}{}
		panic(ErrLockTimeout)
		return
	}
}

func (r *RedisLock) Unlock() {
	value := drivers.GetRedis().Get(r.dc.Key).Val()
	if r.value == value {
		drivers.GetRedis().Del(r.dc.Key)
	}
}

func NewRedisLock (dc DistributedLockConfig) (LockDriver, error) {
	err := drivers.OpenRedis(dc.RedisConfig)
	if err != nil {
		return nil, err
	}
	return &RedisLock{
		dc:          dc,
		lockChannel: make(chan bool),
		exitChannel: make(chan struct{}),
		value: utils.RandNum(),
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
