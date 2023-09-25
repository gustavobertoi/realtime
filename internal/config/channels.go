package config

import (
	"context"

	"github.com/open-source-cloud/realtime/internal/channels"
	redis_adapter "github.com/open-source-cloud/realtime/pkg/redis"
)

const (
	redisAdapter = "REDIS"
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

	pubSub := c.yamlConfig.PubSub
	var channelAdapter *channelAdapter
	ca, err := createChannelAdapter(pubSub, channelAdapter)
	if err != nil {
		return err
	}
	channelAdapter = ca

	channelsMap := c.yamlConfig.Channels
	for _, channelDTO := range channelsMap {
		channel, err := channels.NewChannel(&channels.CreateChannelDTO{
			ID:                      channelDTO.ID,
			Name:                    channelDTO.Name,
			Type:                    channelDTO.Type,
			MaxOfChannelConnections: channelDTO.Config.MaxOfConnections,
		}, channelAdapter.consumer, channelAdapter.producer)
		if err != nil {
			return err
		}
		c.channelStore.Set(channel.ID, channel)
	}

	return nil
}

func createChannelAdapter(pubSub *YamlPubSubDTO, channelAdapter *channelAdapter) (*channelAdapter, error) {
	switch pubSub.Driver {
	case redisAdapter:
		if pubSub.Redis == nil {
			return nil, errRedisPubSubAdapterNotDefined
		}
		channelAdapter, err := createRedisChannelAdapter(pubSub)
		if err != nil {
			return nil, err
		}
		return channelAdapter, nil
	}
	return nil, errDriverNotSupported
}

func createRedisChannelAdapter(pubsub *YamlPubSubDTO) (*channelAdapter, error) {
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
