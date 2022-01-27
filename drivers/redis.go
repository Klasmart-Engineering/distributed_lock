package drivers

import (
	"fmt"
	rd "github.com/go-redis/redis"
	"sync"
)

var (
	redis    *rd.Client
	curConfig RedisConfig
	_openOnce sync.Once
)

type RedisConfig struct {
	Host string
	Port int
	Password string
}

func checkConfig(config1, config2 RedisConfig) bool{
	if config1.Host != config2.Host {
		return false
	}
	if config1.Port != config2.Port {
		return false
	}
	if config1.Password != config2.Password {
		return false
	}
	return true
}

func OpenRedis(config RedisConfig) error{
	_openOnce.Do(func() {
		//若已连接，且配置相同，则直接返回
		//if redis != nil && checkConfig(config, curConfig){
		//	return nil
		//}

		//连接redis
		redis = rd.NewClient(&rd.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
			Password: config.Password, // no password set
			DB:       0,                     // use default DB
		})
		//测试Redis是否连接成功
		//_, err := redis.Ping().Result()
		//if err != nil {
		//	redis = nil
		//	return err
		//}
		curConfig = config
	})

	return nil
}

//Close关闭数据库
func CloseRedis() {
	redis.Close()
}

//GetRedis获得Redis句柄
func GetRedis() *rd.Client {
	return redis
}
