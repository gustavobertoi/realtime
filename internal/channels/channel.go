package channels

import "fmt"

type ChannelConfig struct {
	MaxOfConnections int `json:"max_of_connections"`
}

type Channel struct {
	adapter ClientAdapter
	clients map[string]*Client

	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Config *ChannelConfig `json:"config"`
}

func NewChannel(id string, name string, maxOfConnections int) (*Channel, error) {
	c := &Channel{
		ID:   id,
		Name: name,
		Config: &ChannelConfig{
			MaxOfConnections: maxOfConnections,
		},
		clients: make(map[string]*Client),
		adapter: nil,
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

func (c *Channel) CreateClient(userAgent string, ip string) (*Client, error) {
	if (len(c.clients)) >= c.Config.MaxOfConnections {
		return nil, errMaxLimitOfConnections
	}
	client := NewClient(userAgent, ip, c.ID)
	c.clients[client.ID] = client
	return client, nil
}

func (c *Channel) DeleteClient(clientID string) {
	if c.HasClient(clientID) {
		delete(c.clients, clientID)
	}
}

func (c *Channel) HasClient(clientID string) bool {
	return c.clients[clientID] != nil
}

func (c *Channel) CountOfClients() int {
	return len(c.clients)
}

func (c *Channel) BroadcastAll(m *Message) map[string]error {
	var errMap map[string]error = make(map[string]error)
	for _, client := range c.clients {
		if m.ClientID != client.ID {
			err := client.Send(m)
			if err != nil {
				errMap[client.ID] = err
			}
		}
	}
	return errMap
}

func (c *Channel) SetAdapter(a ClientAdapter) {
	c.adapter = a
}

func (c *Channel) Subscribe(client *Client) error {
	if c.adapter == nil {
		return fmt.Errorf("channel %s does not have client adapter setted", c.ID)
	}
	err := c.adapter.Subscribe(client)
	if err != nil {
		return err
	}
	return nil
}
