package redis

import (
	"context"
	"encoding/json"
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
func (c *Cache) GetV(ctx context.Context, key string, value any) error {
	v, err := c.conn.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(v), value)
	if err != nil {
		return err
	}
	return nil
}

// Set implements cache.Cache.
func (c *Cache) Set(ctx context.Context, key string, value any) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.conn.Set(ctx, key, v, 0).Err()
}

// SetWithTimeout implements cache.Cache.
func (c *Cache) SetWithTimeout(ctx context.Context, key string, value any, timeout time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.conn.Set(ctx, key, v, timeout).Err()
}
