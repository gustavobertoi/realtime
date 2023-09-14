package channels

type ClientAdapter interface {
	Subscribe(client *Client) error
}

type MessageAdapter interface {
	Send(message *Message) error
}

type MessageStore interface {
	Get(id string) (*Message, error)
	Has(id string) bool
	Put(message *Message) error
	Delete(message *Message) error
}
