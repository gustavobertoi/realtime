package channels

import (
	"errors"
	"fmt"
	"log"
)

var (
	errMaxLimitOfConnections        = errors.New("max limit of connections for this client")
	errInvalidMaxLimitOfConnections = errors.New("invalid max limit of connections for this channel")
)

type ChannelConfig struct {
	MaxOfConnections int `json:"max_of_connections"`
}

type Channel struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Config  *ChannelConfig `json:"config"`
	clients map[string]*Client
	adapter ChannelMessageAdapter
}

func NewChannel(id string, name string, maxOfConnections int) (*Channel, error) {
	ch := &Channel{
		ID:   id,
		Name: name,
		Config: &ChannelConfig{
			MaxOfConnections: maxOfConnections,
		},
		clients: make(map[string]*Client),
		adapter: nil,
	}
	err := ch.validate()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (ch *Channel) validate() error {
	// load from a env/yaml
	const maxOfConnectionsPerChannel = 100
	if ch.Config.MaxOfConnections > maxOfConnectionsPerChannel {
		return errInvalidMaxLimitOfConnections
	}
	return nil
}

func (ch *Channel) SetAdapter(adapter ChannelMessageAdapter) {
	ch.adapter = adapter
}

func (ch *Channel) CreateClient(userAgent string, ip string) (*Client, error) {
	if (len(ch.clients)) >= ch.Config.MaxOfConnections {
		return nil, errMaxLimitOfConnections
	}
	client := NewClient(userAgent, ip, ch.ID)
	ch.clients[client.ID] = client
	return client, nil
}

func (ch *Channel) DeleteClient(clientID string) {
	if ch.HasClient(clientID) {
		delete(ch.clients, clientID)
	}
}

func (ch *Channel) HasClient(clientID string) bool {
	return ch.clients[clientID] != nil
}

func (ch *Channel) CountOfClients() int {
	return len(ch.clients)
}

func (ch *Channel) BroadcastAll(m *Message) {
	for _, client := range ch.clients {
		if m.ClientID != client.ID {
			err := client.Send(m)
			if err != nil {
				log.Printf("Failed to send message to client %s: %v", client.ID, err)
			}
		}
	}
}

func (ch *Channel) Subscribe(clientID string, msgChannel chan *Message) error {
	if ch.adapter == nil {
		return fmt.Errorf("channel %s does not contains a subscribe adapter method defined", ch.ID)
	}
	return ch.adapter.Subscribe(ch.ID, clientID, msgChannel)
}

func (ch *Channel) DeleteMessage(messageID string) error {
	if ch.adapter == nil {
		return fmt.Errorf("channel %s does not contains a delete adapter method defined", ch.ID)
	}
	return ch.adapter.DeleteMessage(messageID)
}
