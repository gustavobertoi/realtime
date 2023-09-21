package redis_adapter

import (
	"github.com/go-redis/redis/v8"
	"github.com/open-source-cloud/realtime/internal/channels"
)

func (ra *RedisAdapter) Send(message *channels.Message) error {
	topic := message.ChannelID
	key := message.ID
	value, err := message.MessageToJSON()
	if err != nil {
		return err
	}
	streamArgs := &redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{key: value},
	}
	_, err = ra.client.XAdd(ra.ctx, streamArgs).Result()
	return err
}
