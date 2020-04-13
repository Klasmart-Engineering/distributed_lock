package distributed_lock

import (
	"calmisland/distributed_lock/drivers"
	"context"
	"errors"
	"time"
)

const (
	LockDelay = 5
)

var (
	ErrUnknownDistributedLockDriver = errors.New("Unknown distributed lock driver")
	ErrLockTimeout = errors.New("Lock timeout")
)

type LockDriver  interface {
	Lock()
	Unlock()
}

type RedisLock struct {
	dc DistributedLockConfig
}

func (r *RedisLock) Lock() {
	//尝试等待5s
	for i := 0; i < r.dc.RetryLockDuration * 100; i++ {
		ret, _ := drivers.GetRedis().SetNX(r.dc.Key, "1", r.dc.Timeout).Result()
		if ret {
			return
		}

		time.Sleep(time.Duration(10) * time.Millisecond)
	}
}

func (r *RedisLock) Unlock() {
	drivers.GetRedis().Del(r.dc.Key)
}

func NewRedisLock(dc DistributedLockConfig) (LockDriver , error){
	err := drivers.OpenRedis(dc.RedisConfig)
	if err != nil{
		return nil, err
	}
	return &RedisLock{
		dc: dc,
	}, nil
}


type DistributedLockConfig struct {
	Driver string
	RetryLockDuration int
	Key string
	Timeout time.Duration
	Ctx context.Context
	RedisConfig drivers.RedisConfig
}

func NewDistributedLock(dc DistributedLockConfig) (LockDriver , error) {
	if dc.RetryLockDuration <= 0 {
		dc.RetryLockDuration = LockDelay
	}
	if dc.Driver == "redis" {
		//打开redis
		return NewRedisLock(dc)
	}
	return nil, ErrUnknownDistributedLockDriver
}
