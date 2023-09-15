package channels

type ConsumerAdapter interface {
	Subscribe(client *Client) error
}

type ProducerAdapter interface {
	Send(message *Message) error
}

type ClientStore interface {
	Count() int
	All() []*Client
	Get(id string) (*Client, error)
	Has(id string) bool
	Put(client *Client) error
	Delete(client *Client) error
}