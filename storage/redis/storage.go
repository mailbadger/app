package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

// Storage is the interface that needs to be implemented for caching capabilities.
type Storage interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte, duration time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Expire(key string, seconds time.Duration) error
}

// NewRedisStore creates new redis store
func NewRedisStore() (Storage, error) {
	r, err := newRedisClient()
	if err != nil {
		return nil, err
	}

	return &redisStore{r}, nil
}

func GenCacheKey(prefix string, key string) string {
	return prefix + key
}

// newRedisClient creates new redis client
func newRedisClient() (*redis.Client, error) {
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
