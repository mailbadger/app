package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RedisStore struct {
	client *redis.Client
}

func (rs *RedisStore) Get(ctx context.Context, key string) ([]byte, error) {
	return rs.client.WithContext(ctx).Get(key).Bytes()
}

func (rs *RedisStore) Set(ctx context.Context, key string, content []byte, duration time.Duration) error {
	return rs.client.WithContext(ctx).Set(key, content, duration).Err()
}

func (rs *RedisStore) Delete(ctx context.Context, key string) error {
	return rs.client.WithContext(ctx).Del(key).Err()
}

func (rs *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	r, err := rs.client.WithContext(ctx).Exists(key).Result()
	return r != 0, err
}

func (rs *RedisStore) Expire(ctx context.Context, key string, duration time.Duration) error {
	return rs.client.WithContext(ctx).Expire(key, duration).Err()
}
