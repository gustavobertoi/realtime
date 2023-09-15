package channels

import (
	"context"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

type RedisAdapter struct {
	ProducerAdapter
	ConsumerAdapter
	client *redis.Client
}

func NewRedisAdapter(options redis.Options) *RedisAdapter {
	return &RedisAdapter{
		client: redis.NewClient(&options),
	}
}

func (a *RedisAdapter) Send(message *Message) error {
	logger := log.WithFields(log.Fields{
		"message_id": message.ID,
		"channel_id": message.ChannelID,
		"client_id":  message.ClientID,
		"context":    "redis_adapter.go",
	})
	msg, err := message.ToJSON()
	if err != nil {
		return err
	}
	a.client.Publish(ctx, message.ChannelID, msg)
	logger.Printf("message has been published to redis adapter")
	return nil
}

func (a *RedisAdapter) Subscribe(client *Client) error {
	go func() {
		logger := log.WithFields(log.Fields{
			"channel_id": client.ChannelID,
			"client_id":  client.ID,
			"context":    "redis_adapter.go",
		})
		logger.Print("starting subscription to redis adapter")
		redisClient := a.client
		ch := client.GetChan()
		pubsub := redisClient.Subscribe(ctx, client.ChannelID)
		defer pubsub.Close()
		for {
			redisMsg, err := pubsub.ReceiveMessage(ctx)
			logger.Print("processing message from redis adapter")
			if err != nil {
				logger.Errorf("error reading message from redis adapter, details: %v", err)
				continue
			}
			msg, err := FromJSON(string(redisMsg.Payload))
			if err != nil {
				logger.Errorf("error deserializing message from redis adapter, details: %v", err)
				continue
			}
			ch <- msg
		}
	}()
	return nil
}
