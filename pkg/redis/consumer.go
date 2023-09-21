package redis_adapter

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func (ra *RedisAdapter) Subscribe(topic string, callback func(value interface{})) error {
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
					callback(value)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}
