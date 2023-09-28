package config

import (
	"context"

	"github.com/open-source-cloud/realtime/internal/channels"
	memory_adapter "github.com/open-source-cloud/realtime/pkg/pubsub/memory"
	redis_adapter "github.com/open-source-cloud/realtime/pkg/pubsub/redis"
)

// PubSub adapters
const (
	memoryDriver = "MEMORY"
	redisDriver  = "REDIS"
)

type channelAdapter struct {
	consumer channels.ConsumerAdapter
	producer channels.ProducerAdapter
}

func (c *Config) GetChannelByID(refID string) (*channels.Channel, error) {
	channel, err := c.channelStore.Get(refID)
	if err != nil {
		return nil, err
	}
	return channel.(*channels.Channel), nil
}

func (c *Config) CreateChannelsFromConfig() error {
	if c.rootConfig == nil {
		return errYamlConfigNotDefined
	}
	ps := c.rootConfig.PubSub
	ca, err := createChannelAdapter(ps)
	if err != nil {
		return err
	}
	if len(c.rootConfig.Channels) >= 1 {
		for _, dto := range c.rootConfig.Channels {
			channel, err := channels.NewChannel(&channels.CreateChannelDTO{
				ID:                      dto.ID,
				Name:                    dto.Name,
				Type:                    dto.Type,
				MaxOfChannelConnections: dto.Config.MaxOfConnections,
			}, ca.consumer, ca.producer)
			if err != nil {
				return err
			}
			c.channelStore.Set(channel.ID, channel)
		}
	}
	return nil
}

func createChannelAdapter(pubsub *PubSubDTO) (*channelAdapter, error) {
	switch pubsub.Driver {
	case memoryDriver:
		return createMemoryChannelAdapter(pubsub)
	case redisDriver:
		return createRedisChannelAdapter(pubsub)
	}
	return nil, errDriverNotSupported
}

func createMemoryChannelAdapter(pubsub *PubSubDTO) (*channelAdapter, error) {
	memoryAdapter := memory_adapter.NewMemmoryAdapter()
	return &channelAdapter{
		consumer: memoryAdapter,
		producer: memoryAdapter,
	}, nil
}

func createRedisChannelAdapter(pubsub *PubSubDTO) (*channelAdapter, error) {
	if pubsub.Redis == nil {
		return nil, errRedisPubSubAdapterNotDefined
	}
	redisAdapter, err := redis_adapter.NewRedisAdapter(context.Background(), &redis_adapter.RedisConfig{
		URL: pubsub.Redis.URL,
	})
	if err != nil {
		return nil, err
	}
	return &channelAdapter{
		consumer: redisAdapter,
		producer: redisAdapter,
	}, nil
}
