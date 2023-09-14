package channels

import (
	"fmt"

	"github.com/open-source-cloud/realtime/pkg/uuid"
)

type Client struct {
	adapter MessageAdapter
	store   MessageStore
	ch      chan *Message

	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
}

func NewClient(userAgent string, ipAddress string, channelID string) *Client {
	client := &Client{
		ID:        uuid.NewUUID(),
		ChannelID: channelID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		adapter:   nil,
		store:     nil,
		ch:        make(chan *Message),
	}
	return client
}

func (c *Client) SetAdapter(adapter MessageAdapter) {
	c.adapter = adapter
}

func (c *Client) SetStore(store MessageStore) {
	c.store = store
}

func (c *Client) Send(msg *Message) error {
	if c.adapter == nil {
		return fmt.Errorf("client %s from channel %s does not contains message adapter", c.ID, c.ChannelID)
	}
	if c.store == nil {
		return fmt.Errorf("client %s from channel %s does not contains message store", c.ID, c.ChannelID)
	}
	if c.store.Has(msg.ID) {
		return errMessageAlreadyPublished
	}
	err := c.adapter.Send(msg)
	if err != nil {
		return err
	}
	c.store.Put(msg)
	return nil
}

func (c *Client) GetInternalChannel() chan *Message {
	return c.ch
}

func (c *Client) ProcessMessage(m *Message) {
	m.SetAsPublished()
}

func (c *Client) DeleteMessage(m *Message) error {
	if c.store == nil {
		return fmt.Errorf("client %s from channel %s does not contains message store", c.ID, c.ChannelID)
	}
	return c.store.Delete(m)
}
