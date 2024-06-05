package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/gofreego/goutils/cache/redis"
)

const (
	REDIS = "redis"
)

type Cache interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (string, error)
	GetV(ctx context.Context, key string, value any) error
	SetWithTimeout(ctx context.Context, key string, value any, timeout time.Duration) error
}

type Config struct {
	Name  string
	Redis redis.Config
}

func NewCache(ctx context.Context, conf *Config) Cache {
	switch conf.Name {
	case REDIS:
		return redis.NewCache(ctx, &conf.Redis)
	}
	panic(fmt.Sprintf("invalid cache name, provided %s ,expected : %s", conf.Name, REDIS))
}
