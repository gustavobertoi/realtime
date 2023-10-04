package channels

type (
	CreateChannelDTO struct {
		ID                      string `json:"id" yaml:"id"`
		Type                    string `json:"type" yaml:"type"`
		MaxOfChannelConnections int    `json:"max_of_channel_connections" yaml:"max_of_channel_connections"`
	}
	CreateClientDTO struct {
		ID        string
		ChannelID string
		UserAgent string
		IPAddress string
	}
	SendServerMessageDTO struct {
		Payload string `json:"payload"`
	}
)
