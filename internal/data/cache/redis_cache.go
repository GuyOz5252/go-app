package cache

import (
	"context"
	"time"

	"github.com/GuyOz5252/go-app/internal/core"
)

type RedisCache struct {
}

func NewRedisCache() core.Cache {
	return &RedisCache{}
}

func (c *RedisCache) SetKey(ctx context.Context, key string, ttl time.Duration) error {
	panic("unimplemented")
}

func (c *RedisCache) DeleteKey(ctx context.Context, key string) error {
	panic("unimplemented")
}

func (c *RedisCache) KeyExists(ctx context.Context, key string) (bool, error) {
	panic("unimplemented")
}

func (c *RedisCache) AddToSet(ctx context.Context, key, value string) error {
	panic("unimplemented")
}

func (c *RedisCache) RemoveFromSet(ctx context.Context, key, value string) error {
	panic("unimplemented")
}

func (c *RedisCache) GetSet(ctx context.Context, key string) ([]string, error) {
	panic("unimplemented")
}
