package channels

type ChannelMessageAdapter interface {
	Subscribe(channelID string, clientID string, msgChannel chan *Message) error
	DeleteMessage(messageID string) error
}

type ClientMessageAdapter interface {
	Send(msg *Message) error
}
