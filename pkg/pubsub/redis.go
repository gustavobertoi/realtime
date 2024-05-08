package pubsub

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gustavobertoi/realtime/internal/channels"
	"github.com/gustavobertoi/realtime/internal/dtos"
)

type RedisAdapter struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisAdapter(ctx context.Context, c *dtos.RedisConfig) (*PubSub, error) {
	if c == nil {
		return nil, NewPubSubError(RedisDriver, "redis config is nil")
	}
	opt, err := parseRedisURL(c.URL)
	if err != nil {
		return nil, NewPubSubError(RedisDriver, fmt.Sprintf("error parsing redis URL: %s", err.Error()))
	}
	ra := &RedisAdapter{
		client: redis.NewClient(opt),
		ctx:    ctx,
	}
	return &PubSub{
		Driver:   RedisDriver,
		Consumer: ra,
		Producer: ra,
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

func (ra *RedisAdapter) Send(message *channels.Message) error {
	topic := message.ChannelID
	key := message.ID
	value, err := message.MessageToJSON()
	if err != nil {
		return NewPubSubError(RedisDriver, fmt.Sprintf("error serializing message to json: %s", err.Error()))
	}
	streamArgs := &redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{key: value},
	}
	return ra.client.XAdd(ra.ctx, streamArgs).Err()
}

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
				log.Printf("error reading topic %s from redis stream adapter, err: %v", topic, err)
				continue
			}
			stream := streams[0]
			for _, xMsg := range stream.Messages {
				for _, value := range xMsg.Values {
					rawMsg := value.(string)
					msg, err := channels.MessageFromJSON(rawMsg)
					if err != nil {
						log.Printf("error deserializing message from json, topic %s, err: %v", topic, err)
						continue
					}
					client.MessageChan() <- msg
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}
