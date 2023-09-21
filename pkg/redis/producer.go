package redis_adapter

import (
	"github.com/go-redis/redis/v8"
)

func (ra *RedisAdapter) Send(topic string, key string, value string) error {
	streamArgs := &redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{key: value},
	}
	_, err := ra.client.XAdd(ra.ctx, streamArgs).Result()
	return err
}
