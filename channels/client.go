package channels

import "github.com/open-source-cloud/realtime/pkg/uuid"

type Client struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
	msgChan   chan *Message
}

func NewClient(data *CreateClientDTO) *Client {
	if data.ID == "" {
		data.ID = uuid.NewUUID()
	}
	client := &Client{
		ID:        data.ID,
		ChannelID: data.ChannelID,
		UserAgent: data.UserAgent,
		IPAddress: data.IPAddress,
		msgChan:   make(chan *Message),
	}
	return client
}

func (c *Client) MessageChan() chan *Message {
	return c.msgChan
}
