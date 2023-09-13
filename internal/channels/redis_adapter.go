package channels

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisAdapter struct {
	ChannelMessageAdapter
	ClientMessageAdapter
	rdb       *redis.Client
	mapOfMsgs map[string]string
}

func NewRedisAdapter(options redis.Options) *RedisAdapter {
	return &RedisAdapter{
		rdb:       redis.NewClient(&options),
		mapOfMsgs: make(map[string]string),
	}
}

func (adapter *RedisAdapter) hasMsg(msgID string) bool {
	return adapter.mapOfMsgs[msgID] != ""
}

func (adapter *RedisAdapter) setMsg(msgID string) {
	adapter.mapOfMsgs[msgID] = msgID
}

func (adapter *RedisAdapter) Send(m *Message) error {
	if !adapter.hasMsg(m.ID) {
		redisClient := adapter.rdb
		msg, err := m.ToJSON()
		if err != nil {
			return err
		}
		redisClient.Publish(ctx, m.ChannelID, msg)
		adapter.setMsg(m.ID)
	}
	return nil
}

func (adapter *RedisAdapter) Subscribe(channelID string, clientID string, msgChannel chan *Message) error {
	redisClient := adapter.rdb
	go func() {
		pubsub := redisClient.Subscribe(ctx, channelID)
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
			msgChannel <- m
		}
	}()
	return nil
}

func (adapter *RedisAdapter) DeleteMessage(messageID string) error {
	if adapter.hasMsg(messageID) {
		log.Printf("deleting messing from redis adapter %s", messageID)
		delete(adapter.mapOfMsgs, messageID)
	}
	return nil
}
