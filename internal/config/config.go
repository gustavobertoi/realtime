package config

import (
	"context"
	"errors"

	"github.com/open-source-cloud/realtime/internal/channels"
	redis_adapter "github.com/open-source-cloud/realtime/pkg/redis"
)

var errChannelNotFound = errors.New("channel not found")

var redisAdapter = redis_adapter.NewRedisAdapter(context.Background(), &redis_adapter.RedisConfig{
	URL: "redis://default:realtime@localhost:6379",
})
var channelsRedisAdapter = channels.NewChannelsRedisAdapter(redisAdapter)

type Config struct {
	Port        int
	channelsMap map[string]*channels.Channel
}

func NewConfig() *Config {
	c := &Config{
		Port:        8080,
		channelsMap: make(map[string]*channels.Channel),
	}
	c.loadChannelsFromYaml()
	return c
}

// TODO: Refactor this function to load from yaml
func (c *Config) loadChannelsFromYaml() {
	var eventsChannel, err = channels.NewChannel("742fc7fe-1527-4184-8945-10b30bf01347", "events", 2, channelsRedisAdapter)
	if (err) != nil {
		panic(err)
	}
	c.channelsMap[eventsChannel.ID] = eventsChannel
}

func (c *Config) GetChannelByID(id string) (*channels.Channel, error) {
	ch := c.channelsMap[id]
	if ch == nil {
		return nil, errChannelNotFound
	}
	return ch, nil
}
