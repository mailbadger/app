package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type redisStore struct {
	client *redis.Client
}

func NewRedisClient() (*redis.Client, error) {
	opts := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DB:   0,
	}
	if os.Getenv("REDIS_PASS") != "" {
		opts.Password = os.Getenv("REDIS_PASS")
	}

	client := redis.NewClient(opts)
	_, err := client.Ping().Result()
	if err != nil {
		// try without password
		opts.Password = ""
		client = redis.NewClient(opts)
		_, err = client.Ping().Result()
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func NewRedisStore() (Storage, error) {
	r, err := NewRedisClient()
	if err != nil {
		return nil, err
	}

	return &redisStore{r}, nil
}

func (rs *redisStore) Get(key string) ([]byte, error) {
	return rs.client.Get(key).Bytes()
}

func (rs *redisStore) Set(key string, content []byte, duration time.Duration) error {
	s := rs.client.Set(key, content, duration)
	return s.Err()
}

func (rs *redisStore) Delete(key string) error {
	return rs.client.Del(key).Err()
}

func (rs *redisStore) Exists(key string) bool {
	r, err := rs.client.Exists(key).Result()
	return r != 0 && err == nil
}

func (rs *redisStore) Expire(key string, duration time.Duration) error {
	return rs.client.Expire(key, duration).Err()
}
