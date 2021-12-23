package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"jwtToken/cfg"
	"log"
)

var redisClient *redis.Client

// 创建 redis 客户端
func createClient() *redis.Client {

	cfg := cfg.Cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	})

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

func InitRedis() {
	redisClient = createClient()
}

func RedisClient() *redis.Client {
	return redisClient
}
