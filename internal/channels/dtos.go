package channels

type (
	CreateChannelDTO struct {
		ID                      string
		Name                    string
		MaxOfChannelConnections int
		Type                    string
	}
	CreateClientDTO struct {
		ID        string
		ChannelID string
		UserAgent string
		IPAddress string
	}
)
