package channels

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisAdapter struct {
	ClientAdapter
	MessageAdapter
	client *redis.Client
}

func NewRedisAdapter(options redis.Options) *RedisAdapter {
	return &RedisAdapter{
		client: redis.NewClient(&options),
	}
}

func (a *RedisAdapter) Send(message *Message) error {
	msg, err := message.ToJSON()
	if err != nil {
		return err
	}
	a.client.Publish(ctx, message.ChannelID, msg)
	return nil
}

func (a *RedisAdapter) Subscribe(client *Client) error {
	go func() {
		redisClient := a.client
		channelID := client.ChannelID
		clientID := client.ID
		ch := client.GetInternalChannel()
		pubsub := redisClient.Subscribe(ctx, client.ChannelID)
		defer pubsub.Close()
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Printf("error reading from redis pubsub (channel %s - client %s): %v", channelID, clientID, err)
				return
			}
			m, err := FromJSON(string(msg.Payload))
			if err != nil {
				fmt.Printf("error parsing message from redis pubsub (channel %s - client %s): %v", channelID, clientID, err)
				return
			}
			ch <- m
		}
	}()
	return nil
}
