package channels

type (
	CreateChannelDTO struct {
		ID                      string `json:"id"`
		Type                    string `json:"type"`
		MaxOfChannelConnections int    `json:"max_of_connections"`
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
