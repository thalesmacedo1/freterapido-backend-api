package redis

import (
	"context"
	"strings"
	"time"
)

type Cache interface {
	SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetCache(ctx context.Context, key string) (string, error)
	DeleteCache(ctx context.Context, key string) error
	ExistsCache(ctx context.Context, key string) (bool, error)
	SetCacheJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetCacheJSON(ctx context.Context, key string, dest interface{}) error
}

type RedisCache struct {
	client *RedisClient
}

func NewRedisCache(client *RedisClient) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.SetCache(ctx, key, value, expiration)
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return rc.client.GetCache(ctx, key)
}

func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.DeleteCache(ctx, key)
}

func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	return rc.client.ExistsCache(ctx, key)
}

func (rc *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.SetCacheJSON(ctx, key, value, expiration)
}

func (rc *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	return rc.client.GetCacheJSON(ctx, key, dest)
}

func GenerateCacheKey(parts ...string) string {
	return strings.Join(parts, ":")
}
