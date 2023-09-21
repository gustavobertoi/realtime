package redis_adapter

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/open-source-cloud/realtime/internal/channels"
)

type RedisConfig struct {
	Addr     string
	Password string
}

type RedisAdapter struct {
	channels.ProducerAdapter
	channels.ConsumerAdapter
	client *redis.Client
	ctx    context.Context
}

func NewRedisAdapter(ctx context.Context, c *RedisConfig) *RedisAdapter {
	return &RedisAdapter{
		client: redis.NewClient(&redis.Options{
			Addr:     c.Addr,
			Password: c.Password,
		}),
		ctx: ctx,
	}
}
