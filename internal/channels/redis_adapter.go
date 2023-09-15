package channels

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/open-source-cloud/realtime/pkg/log"
	"github.com/sirupsen/logrus"
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
	logger := log.CreateWithContext("redis_adapter.go", logrus.Fields{
		"message_id": message.ID,
		"channel_id": message.ChannelID,
		"client_id":  message.ClientID,
	})
	msg, err := message.ToJSON()
	if err != nil {
		return err
	}
	streamArgs := &redis.XAddArgs{
		Stream: message.ChannelID,
		Values: map[string]interface{}{message.ID: msg},
	}
	if _, err := a.client.XAdd(ctx, streamArgs).Result(); err != nil {
		return err
	}
	logger.Printf("message has been published to redis stream adapter")
	return nil
}

func (a *RedisAdapter) Subscribe(client *Client) error {
	go func() {
		logger := log.CreateWithContext("redis_adapter.go", logrus.Fields{
			"channel_id": client.ChannelID,
			"client_id":  client.ID,
		})
		logger.Print("starting subscription to redis stream adapter")
		ch := client.GetChan()
		redisClient := a.client
		for {
			streamArgs := redis.XReadArgs{
				Streams: []string{client.ChannelID, "$"},
				Block:   0,
				Count:   0,
			}
			streams, err := redisClient.XRead(ctx, &streamArgs).Result()
			if err != nil {
				logger.Errorf("could not read from the stream: %v", err)
				continue
			}
			stream := streams[0]
			for _, xMsg := range stream.Messages {
				for _, value := range xMsg.Values {
					logger.Printf("processing msg from redis stream, casting to string")
					str, ok := value.(string)
					if !ok {
						logger.Errorf("failed to cast to value from redis stream to string")
						continue
					}
					msg, err := FromJSON(str)
					if err != nil {
						logger.Errorf("error deserializing message from redis adapter, details: %v", err)
						continue
					}
					logger.Printf("message casted writing to internal client channel")
					ch <- msg
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
	return nil
}
