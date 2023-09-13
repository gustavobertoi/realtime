package channels

import (
	"fmt"

	"github.com/open-source-cloud/realtime/pkg/uuid"
)

type Client struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
	adapter   ClientMessageAdapter
}

func NewClient(userAgent string, ipAddress string, channelID string) *Client {
	client := &Client{
		ID:        uuid.NewUUID(),
		ChannelID: channelID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		adapter:   nil,
	}
	return client
}

func (c *Client) SetAdapter(adapter ClientMessageAdapter) {
	c.adapter = adapter
}

func (c *Client) Send(msg *Message) error {
	if c.adapter == nil {
		return fmt.Errorf("client %s of channel %s does not contains send adapter", c.ID, c.ChannelID)
	}
	return c.adapter.Send(msg)
}
