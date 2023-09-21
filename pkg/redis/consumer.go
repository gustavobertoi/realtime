package redis_adapter

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/open-source-cloud/realtime/internal/channels"
)

func (ra *RedisAdapter) Subscribe(client *channels.Client) error {
	topic := client.ChannelID
	go func() {
		for {
			streamArgs := redis.XReadArgs{
				Streams: []string{topic, "$"},
				Block:   0,
				Count:   0,
			}
			streams, err := ra.client.XRead(ra.ctx, &streamArgs).Result()
			if err != nil {
				log.Panicf("error reading topic %s from redis stream adapter, err: %v", topic, err)
			}
			stream := streams[0]
			for _, xMsg := range stream.Messages {
				for _, value := range xMsg.Values {
					rawMsg := value.(string)
					msg, err := channels.MessageFromJSON(rawMsg)
					if err != nil {
						log.Panicf("error deserializing message from json, topic %s, err: %v", topic, err)
					}
					client.MessageChan() <- msg
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}
