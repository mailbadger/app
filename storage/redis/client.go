package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type redisStore struct {
	client *redis.Client
}

func (rs *redisStore) Get(key string) ([]byte, error) {
	return rs.client.Get(key).Bytes()
}

func (rs *redisStore) Set(key string, content []byte, duration time.Duration) error {
	return rs.client.Set(key, content, duration).Err()
}

func (rs *redisStore) Delete(key string) error {
	return rs.client.Del(key).Err()
}

func (rs *redisStore) Exists(key string) (bool, error) {
	r, err := rs.client.Exists(key).Result()
	return r != 0, err
}

func (rs *redisStore) Expire(key string, duration time.Duration) error {
	return rs.client.Expire(key, duration).Err()
}
