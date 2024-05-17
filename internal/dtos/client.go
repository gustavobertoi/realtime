package dtos

type CreateClientDTO struct {
	ID            string
	ChannelID     string
	UserAgent     string
	IPAddress     string
	MaxOfMessages int
}
