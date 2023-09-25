package config

import (
	"context"

	"github.com/open-source-cloud/realtime/internal/channels"
	memory_adapter "github.com/open-source-cloud/realtime/pkg/pubsub/memory"
	redis_adapter "github.com/open-source-cloud/realtime/pkg/pubsub/redis"
)

// PubSub adapters
const (
	memoryAdapter = "MEMORY"
	redisAdapter  = "REDIS"
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

func (c *Config) createChannelsFromConfig() error {
	if c.yamlConfig == nil {
		return errYamlConfigNotDefined
	}

	ps := c.yamlConfig.PubSub
	ca, err := createChannelAdapter(ps)
	if err != nil {
		return err
	}

	channelsMap := c.yamlConfig.Channels
	for _, channelDTO := range channelsMap {
		channel, err := channels.NewChannel(&channels.CreateChannelDTO{
			ID:                      channelDTO.ID,
			Name:                    channelDTO.Name,
			Type:                    channelDTO.Type,
			MaxOfChannelConnections: channelDTO.Config.MaxOfConnections,
		}, ca.consumer, ca.producer)
		if err != nil {
			return err
		}
		c.channelStore.Set(channel.ID, channel)
	}

	return nil
}

func createChannelAdapter(pubsub *YamlPubSubDTO) (*channelAdapter, error) {
	switch pubsub.Driver {
	case memoryAdapter:
		return createMemoryChannelAdapter(pubsub)
	case redisAdapter:
		return createRedisChannelAdapter(pubsub)
	}
	return nil, errDriverNotSupported
}

func createMemoryChannelAdapter(pubsub *YamlPubSubDTO) (*channelAdapter, error) {
	memoryAdapter := memory_adapter.NewMemmoryAdapter()
	return &channelAdapter{
		consumer: memoryAdapter,
		producer: memoryAdapter,
	}, nil
}

func createRedisChannelAdapter(pubsub *YamlPubSubDTO) (*channelAdapter, error) {
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
