package redis_adapter

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/open-source-cloud/realtime/internal/channels"
)

type RedisConfig struct {
	URL string
}

type RedisAdapter struct {
	channels.ProducerAdapter
	channels.ConsumerAdapter
	client *redis.Client
	ctx    context.Context
}

func NewRedisAdapter(ctx context.Context, c *RedisConfig) (*RedisAdapter, error) {
	opt, err := parseRedisURL(c.URL)
	if err != nil {
		return nil, err
	}
	return &RedisAdapter{
		client: redis.NewClient(opt),
		ctx:    ctx,
	}, nil
}

func parseRedisURL(redisURL string) (*redis.Options, error) {
	u, err := url.Parse(redisURL)
	if err != nil {
		return nil, err
	}
	var password string
	if u.User != nil {
		password, _ = u.User.Password()
	}
	var db int
	if len(u.Path) > 1 {
		dbString := strings.TrimPrefix(u.Path, "/")
		db, err = strconv.Atoi(dbString)
		if err != nil {
			return nil, fmt.Errorf("invalid DB number: %v", err)
		}
	}
	return &redis.Options{
		Addr:     u.Host,
		Password: password,
		DB:       db,
	}, nil
}
