package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	conn *redis.Client
}

type Config struct {
	Address  string        // Redis server address, e.g., "localhost:6379"
	Password string        // Password for Redis server, if any
	DB       int           // Redis database to connect to
	PoolSize int           // Maximum number of connections in the pool
	Timeout  time.Duration // Connection timeout duration

}

func NewCache(ctx context.Context, conf *Config) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:        conf.Address,
		Password:    conf.Password,
		DB:          conf.DB,
		PoolSize:    conf.PoolSize,
		DialTimeout: conf.Timeout,
	})

	// Ping the Redis server to check the connection
	err := client.Ping(ctx).Err()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to redis, Err:  %v", err.Error()))
	}

	return &Cache{conn: client}
}

// GetV implements cache.Cache.
func (c *Cache) GetV(key string, value any) error {
	panic("unimplemented")
}

// Get implements cache.Cache.
func (c *Cache) Get(key string) (string, error) {
	panic("unimplemented")
}

// Set implements cache.Cache.
func (c *Cache) Set(key string, value any) error {
	panic("unimplemented")
}

// SetWithTimeout implements cache.Cache.
func (c *Cache) SetWithTimeout(key string, value any, timeout time.Duration) error {
	panic("unimplemented")
}
