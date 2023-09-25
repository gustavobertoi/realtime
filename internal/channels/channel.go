package channels

import "github.com/open-source-cloud/realtime/pkg/store"

const (
	WebSocket        = "WS"
	ServerSentEvents = "SSE"
)

type (
	ChannelConfig struct {
		MaxOfChannelConnections int `json:"max_of_channels_connections"`
	}
	Channel struct {
		ID     string         `json:"id"`
		Name   string         `json:"name"`
		Config *ChannelConfig `json:"config"`
		Type   string         `json:"type"`

		store    *store.MemoryStore
		producer ProducerAdapter
		consumer ConsumerAdapter
	}
)

func NewChannel(dto *CreateChannelDTO, consumer ConsumerAdapter, producer ProducerAdapter) (*Channel, error) {
	c := &Channel{
		ID:   dto.ID,
		Name: dto.Name,
		Config: &ChannelConfig{
			MaxOfChannelConnections: dto.MaxOfChannelConnections,
		},
		Type:     dto.Type,
		store:    store.NewMemoryStore(),
		consumer: consumer,
		producer: producer,
	}
	err := c.validate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Channel) validate() error {
	if c.Config.MaxOfChannelConnections <= 0 {
		return errInvalidMaxLimitOfConnections
	}
	if c.Type != WebSocket && c.Type != ServerSentEvents {
		return errInvalidChannelType
	}
	return nil
}

func (c *Channel) IsMaxOfConnections() bool {
	return c.store.Count() >= c.Config.MaxOfChannelConnections
}

func (c *Channel) BroadcastMessage(m *Message) error {
	return c.producer.Send(m)
}

func (c *Channel) Subscribe(client *Client) error {
	return c.consumer.Subscribe(client)
}

func (c *Channel) DeleteClient(client *Client) {
	c.store.Delete(client.ID)
}

func (c *Channel) AddClient(client *Client) {
	c.store.Set(client.ID, client)
}

func (c *Channel) IsWebSocket() bool {
	return c.Type == WebSocket
}

func (c *Channel) IsSSE() bool {
	return c.Type == ServerSentEvents
}
