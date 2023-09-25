package channels

type (
	CreateChannelDTO struct {
		ID                      string `json:"id"`
		Name                    string `json:"name"`
		MaxOfChannelConnections int    `json:"max_of_connections"`
		Type                    string `json:"type"`
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
