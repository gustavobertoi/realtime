package channels

import "fmt"

type ChannelConfig struct {
	MaxOfConnections int `json:"maxOfConnections"`
}

type Channel struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Config *ChannelConfig `json:"config"`

	consumer    ConsumerAdapter
	clientStore ClientStore
}

func NewChannel(id string, name string, maxOfConnections int) (*Channel, error) {
	c := &Channel{
		ID:   id,
		Name: name,
		Config: &ChannelConfig{
			MaxOfConnections: maxOfConnections,
		},
		consumer:    nil,
		clientStore: nil,
	}
	err := c.validate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Channel) validate() error {
	const maxOfConnectionsPerChannel = 100
	if c.Config.MaxOfConnections > maxOfConnectionsPerChannel {
		return errInvalidMaxLimitOfConnections
	}
	return nil
}

func (c *Channel) BroadcastToAllClients(m *Message) []error {
	var errs []error
	clientStore, err := c.ClientStore()
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	for _, client := range clientStore.All() {
		if m.ClientID != client.ID {
			err := client.Send(m)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func (c *Channel) SetConsumerAdapter(consumer ConsumerAdapter) {
	c.consumer = consumer
}

func (c *Channel) SetClientStore(store ClientStore) {
	c.clientStore = store
}

func (c *Channel) Subscribe(client *Client) error {
	if c.consumer == nil {
		return fmt.Errorf("channel %s does not contains consumer adapter defined", c.ID)
	}
	return c.consumer.Subscribe(client)
}

func (c *Channel) ClientStore() (ClientStore, error) {
	if c.clientStore == nil {
		return nil, fmt.Errorf("channel %s does not contains an client store defined", c.ID)
	}
	return c.clientStore, nil
}
