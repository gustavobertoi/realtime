package channels

type ChannelConfig struct {
	MaxOfChannelConnections int `json:"max_of_channels_connections"`
}

type Channel struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Config   *ChannelConfig `json:"config"`
	Store    *ClientStore
	Consumer ConsumerAdapter
}

func NewChannel(id string, name string, maxOfChannelConnections int, consumer ConsumerAdapter) (*Channel, error) {
	c := &Channel{
		ID:   id,
		Name: name,
		Config: &ChannelConfig{
			MaxOfChannelConnections: maxOfChannelConnections,
		},
		Consumer: consumer,
		Store:    NewClientStore(),
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

func (c *Channel) BroadcastMessage(m *Message) map[string]error {
	errs := make(map[string]error)
	for _, client := range c.Store.All() {
		if m.ClientID != client.ID {
			err := client.Send(m)
			if err != nil {
				errs[client.ID] = err
			}
		}
	}
	return errs
}
