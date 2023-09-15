package channels

import (
	"fmt"
)

type Client struct {
	ID        string `json:"id"`
	ChannelID string `json:"channelId"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`

	producerAdapter ProducerAdapter
	messageStore    MessageStore
	ch              chan *Message
}

func NewClient(clientID string, userAgent string, ipAddress string, channelID string) *Client {
	client := &Client{
		ID:              clientID,
		ChannelID:       channelID,
		UserAgent:       userAgent,
		IPAddress:       ipAddress,
		producerAdapter: nil,
		messageStore:    nil,
		ch:              make(chan *Message),
	}
	return client
}

func (c *Client) SetProducerAdapter(adapter ProducerAdapter) {
	c.producerAdapter = adapter
}

func (c *Client) SetMessageStore(store MessageStore) {
	c.messageStore = store
}

func (c *Client) Send(msg *Message) error {
	if c.producerAdapter == nil {
		return fmt.Errorf("client %s from channel %s does not contains message adapter", c.ID, c.ChannelID)
	}
	if c.messageStore == nil {
		return fmt.Errorf("client %s from channel %s does not contains message store", c.ID, c.ChannelID)
	}
	if c.messageStore.Has(msg.ID) {
		return errMessageAlreadyPublished
	}
	err := c.producerAdapter.Send(msg)
	if err != nil {
		return err
	}
	msg.SetAsPublished()
	c.messageStore.Put(msg)
	return nil
}

func (c *Client) GetChan() chan *Message {
	return c.ch
}

func (c *Client) MessageStore() (MessageStore, error) {
	if c.messageStore == nil {
		return nil, fmt.Errorf("client %s from channel %s does not contains message store", c.ID, c.ChannelID)
	}
	return c.messageStore, nil
}
