package cache

import (
	"context"
	"time"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(address, password string) core.Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) SetKey(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Set(ctx, key, "1", ttl).Err()
}

func (c *RedisCache) DeleteKey(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) KeyExists(ctx context.Context, key string) (bool, error) {
	val, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (c *RedisCache) AddToSet(ctx context.Context, key, value string) error {
	return c.client.SAdd(ctx, key, value).Err()
}

func (c *RedisCache) RemoveFromSet(ctx context.Context, key, value string) error {
	return c.client.SRem(ctx, key, value).Err()
}

func (c *RedisCache) GetSet(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}
