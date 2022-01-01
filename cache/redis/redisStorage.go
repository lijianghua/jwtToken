package redis

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisStorage struct {
	storage *redis.Client
}

func NewRedisStorage(c *redis.Client) *RedisStorage {
	return &RedisStorage{
		storage: c,
	}
}

func (r RedisStorage) Get(key string) (string, error) {
	return r.storage.Get(key).Result()
}

func (r RedisStorage) Set(key, value string, expiresIn time.Duration) error {
	return r.storage.Set(key, value, expiresIn).Err()
}

func (r RedisStorage) Del(key string) (int64, error) {
	return r.storage.Del(key).Result()

}
