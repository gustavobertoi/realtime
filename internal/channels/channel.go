package channels

type ChannelConfig struct {
	MaxOfChannelConnections int `json:"max_of_channels_connections"`
}

type Channel struct {
	producer ProducerAdapter
	consumer ConsumerAdapter

	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Config *ChannelConfig `json:"config"`
	Store  *ClientStore
}

func NewChannel(id string, name string, maxOfChannelConnections int, consumer ConsumerAdapter, producer ProducerAdapter) (*Channel, error) {
	c := &Channel{
		ID:   id,
		Name: name,
		Config: &ChannelConfig{
			MaxOfChannelConnections: maxOfChannelConnections,
		},
		Store:    NewClientStore(),
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
	const maxOfConnectionsPerChannel = 100
	if c.Config.MaxOfChannelConnections > maxOfConnectionsPerChannel {
		return errInvalidMaxLimitOfConnections
	}
	return nil
}

func (c *Channel) IsMaxOfConnections() bool {
	return c.Store.Count() >= c.Config.MaxOfChannelConnections
}

func (c *Channel) BroadcastMessage(m *Message) error {
	return c.producer.Send(m)
}

func (c *Channel) Subscribe(client *Client) error {
	return c.consumer.Subscribe(client)
}
