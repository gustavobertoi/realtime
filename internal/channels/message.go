package channels

import (
	"encoding/json"
	"time"

	"github.com/open-source-cloud/realtime/pkg/uuid"
)

type Message struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channelId"`
	ClientID  string    `json:"clientId"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewMessage(channelID string, clientID string, payload string) *Message {
	id := uuid.NewUUID()
	return &Message{
		ID:        id,
		ChannelID: channelID,
		ClientID:  clientID,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
}

func FromJSON(message string) (*Message, error) {
	m := Message{}
	b := []byte(message)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Message) ToJSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
