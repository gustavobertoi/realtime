package channels

import (
	"log"

	redis_adapter "github.com/open-source-cloud/realtime/pkg/redis"
)

type ChannelsRedisAdapter struct {
	ProducerAdapter
	ConsumerAdapter
	adapter *redis_adapter.RedisAdapter
}

func NewChannelsRedisAdapter(adapter *redis_adapter.RedisAdapter) *ChannelsRedisAdapter {
	return &ChannelsRedisAdapter{
		adapter: adapter,
	}
}

func (cra *ChannelsRedisAdapter) Send(m *Message) error {
	str, err := m.MessageToJSON()
	if err != nil {
		return err
	}
	return cra.adapter.Send(m.ChannelID, m.ID, str)
}

func (cra *ChannelsRedisAdapter) Subscribe(c *Client) error {
	return cra.adapter.Subscribe(c.ChannelID, func(value interface{}) {
		m, err := MessageFromJSON(value.(string))
		if err != nil {
			log.Panicf("error deserializing msg from json, client %s channel %s", c.ID, c.ChannelID)
		}
		ch := c.GetChan()
		ch <- m
	})
}
