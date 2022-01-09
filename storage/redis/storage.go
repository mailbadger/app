package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/mailbadger/app/config"
)

// Store is the interface that needs to be implemented for caching capabilities.
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, content []byte, duration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, seconds time.Duration) error
}

func NewStoreFrom(conf config.Config) (*RedisStore, error) {
	client, err := NewRedisClient(conf.Storage.Redis.Host, conf.Storage.Redis.Port, conf.Storage.Redis.Pass)
	if err != nil {
		return nil, err
	}
	return &RedisStore{client}, nil
}

func NewStore(client *redis.Client) *RedisStore {
	return &RedisStore{client}
}

// NewRedisClient creates new redis client
func NewRedisClient(host, port, pass string) (*redis.Client, error) {
	opts := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		DB:   0,
	}
	if pass != "" {
		opts.Password = pass
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
