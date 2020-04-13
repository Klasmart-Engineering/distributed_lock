package distributed_lock

import (
	"calmisland/distributed_lock/drivers"
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
	Lock(key string, timeout time.Duration) error
	Unlock(key string)
}

type RedisLock struct {
	dc DistributedLockConfig
}

func (r *RedisLock) Lock(key string, timeout time.Duration) error {
	//尝试等待5s
	for i := 0; i < r.dc.RetryLockDuration * 100; i++ {
		ret, err := drivers.GetRedis().SetNX(key, "1", timeout).Result()
		if err != nil {
			return err
		}
		if ret {
			return nil
		}

		time.Sleep(time.Duration(10) * time.Millisecond)
	}

	return ErrLockTimeout
}

func (r *RedisLock) Unlock(key string) {
	drivers.GetRedis().Del(key)
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
