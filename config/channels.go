package config

import (
	"context"

	"github.com/open-source-cloud/realtime/channels"
	"github.com/open-source-cloud/realtime/pkg/pubsub"
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
			c.CreateChannel(dto, ca)
		}
	}
	return nil
}

func (c *Config) CreateChannel(dto *channels.CreateChannelDTO, ca *channelAdapter) (*channels.Channel, error) {
	if ca == nil {
		memoryAdapter, err := createMemoryChannelAdapter()
		if err != nil {
			return nil, err
		}
		ca = memoryAdapter
	}
	channel, err := channels.NewChannel(
		&channels.CreateChannelDTO{
			ID:                      dto.ID,
			Type:                    dto.Type,
			MaxOfChannelConnections: dto.MaxOfChannelConnections,
		},
		ca.consumer,
		ca.producer,
	)
	if err != nil {
		return nil, err
	}
	c.channelStore.Set(channel.ID, channel)
	return channel, nil
}

func createChannelAdapter(ps *PubSub) (*channelAdapter, error) {
	switch ps.Driver {
	case memoryDriver:
		return createMemoryChannelAdapter()
	case redisDriver:
		return createRedisChannelAdapter(ps)
	}
	return nil, errDriverNotSupported
}

func createMemoryChannelAdapter() (*channelAdapter, error) {
	memoryAdapter := pubsub.NewMemmoryAdapter()
	return &channelAdapter{
		consumer: memoryAdapter,
		producer: memoryAdapter,
	}, nil
}

func createRedisChannelAdapter(ps *PubSub) (*channelAdapter, error) {
	if ps.Redis == nil {
		return nil, errRedisPubSubAdapterNotDefined
	}
	redisAdapter, err := pubsub.NewRedisAdapter(context.Background(), &pubsub.RedisConfig{
		URL: ps.Redis.URL,
	})
	if err != nil {
		return nil, err
	}
	return &channelAdapter{
		consumer: redisAdapter,
		producer: redisAdapter,
	}, nil
}
