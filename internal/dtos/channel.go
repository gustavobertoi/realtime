package dtos

type CreateChannelDTO struct {
	ID                      string `json:"id" yaml:"id"`
	Type                    string `json:"type" yaml:"type"`
	MaxOfChannelConnections int    `json:"maxOfChannelConnections" yaml:"max_of_channel_connections"`
}
