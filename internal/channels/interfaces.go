package channels

type ConsumerAdapter interface {
	Subscribe(client *Client) error
}

type ProducerAdapter interface {
	Send(message *Message) error
}
