package drivers

import (
	"fmt"
	rd "github.com/go-redis/redis"
)

var (
	redis    *rd.Client
)

type RedisConfig struct {
	Host string
	Port int
	Password string
}

func OpenRedis(config RedisConfig) error{
	//连接redis
	redis = rd.NewClient(&rd.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password, // no password set
		DB:       0,                     // use default DB
	})
	//测试Redis是否连接成功
	_, err := redis.Ping().Result()
	if err != nil {
		return err
	}
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
