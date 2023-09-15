package config

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/open-source-cloud/realtime/internal/channels"
)

var errChannelNotFound = errors.New("channel not found")

var redisAdapter = channels.NewRedisAdapter(redis.Options{
	Addr:     "localhost:6379",
	Password: "realtime",
})
var clientMemoryStore = channels.NewClientMemoryStore()

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
	ch.SetConsumerAdapter(c.GetChannelConsumerAdapter())
	ch.SetClientStore(clientMemoryStore)
	return ch, nil
}

func (c *Config) GetChannelConsumerAdapter() channels.ConsumerAdapter {
	return redisAdapter
}

func (c *Config) GetClientProducerAdapter() channels.ProducerAdapter {
	return redisAdapter
}
