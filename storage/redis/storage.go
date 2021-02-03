package redis

import "time"

// Storage is the interface that needs to be implemented for caching capabilities.
type Storage interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte, duration time.Duration) error
	Delete(key string) error
	Exists(key string) bool
	Expire(key string, seconds time.Duration) error
}
