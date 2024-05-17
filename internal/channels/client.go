package channels

import (
	"github.com/gustavobertoi/realtime/internal/dtos"
	"github.com/gustavobertoi/realtime/pkg/uuid"
)

type Client struct {
	ID        string `json:"id"`
	ChannelID string `json:"channelId"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`
	msgChan   chan *Message
}

func NewClient(dto *dtos.CreateClientDTO) *Client {
	if dto.ID == "" {
		dto.ID = uuid.NewUUID()
	}
	client := &Client{
		ID:        dto.ID,
		ChannelID: dto.ChannelID,
		UserAgent: dto.UserAgent,
		IPAddress: dto.IPAddress,
		msgChan:   make(chan *Message, dto.MaxOfMessages),
	}
	return client
}

func (c *Client) MessageChan() chan *Message {
	return c.msgChan
}
