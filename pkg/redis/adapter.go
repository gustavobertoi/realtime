package redis_adapter

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	URL string
}

type RedisAdapter struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisAdapter(ctx context.Context, c *RedisConfig) *RedisAdapter {
	return &RedisAdapter{
		client: redis.NewClient(&redis.Options{
			Addr: c.URL,
		}),
		ctx: ctx,
	}
}
