package channels

import (
	"fmt"
)

type Client struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`

	producerAdapter ProducerAdapter
	ch              chan *Message
}

func NewClient(clientID string, userAgent string, ipAddress string, channelID string) *Client {
	client := &Client{
		ID:              clientID,
		ChannelID:       channelID,
		UserAgent:       userAgent,
		IPAddress:       ipAddress,
		producerAdapter: nil,
		ch:              make(chan *Message),
	}
	return client
}

func (c *Client) SetProducerAdapter(adapter ProducerAdapter) {
	c.producerAdapter = adapter
}

func (c *Client) Send(msg *Message) error {
	if c.producerAdapter == nil {
		return fmt.Errorf("client %s from channel %s does not contains message adapter", c.ID, c.ChannelID)
	}
	err := c.producerAdapter.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetChan() chan *Message {
	return c.ch
}
