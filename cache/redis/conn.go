package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"jwtToken/config"
)

var redisClient *redis.Client

// 创建 redis 客户端
func NewClient(cnf *config.Config) error {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cnf.Redis.Host, cnf.Redis.Port),
		Password: "",
		DB:       0,
	})

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := redisClient.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

func RedisClient() *redis.Client {
	return redisClient
}
