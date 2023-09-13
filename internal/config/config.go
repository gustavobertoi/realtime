package config

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/open-source-cloud/realtime/internal/channels"
)

var redisAdapter = channels.NewRedisAdapter(redis.Options{
	Addr:     "localhost:6379",
	Password: "realtime",
})
var errChannelNotFound = errors.New("channel not found")

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
	var eventsChannel, err = channels.NewChannel("742fc7fe-1527-4184-8945-10b30bf01347", "events", 2)
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

func (c *Config) GetChannelAdapter() channels.ChannelMessageAdapter {
	return redisAdapter
}

func (c *Config) GetClientAdapter() channels.ClientMessageAdapter {
	return redisAdapter
}
