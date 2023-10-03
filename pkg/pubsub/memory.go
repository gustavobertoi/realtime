package pubsub

import (
	"github.com/open-source-cloud/realtime/channels"
)

type MemoryAdapter struct {
	channels.ProducerAdapter
	channels.ConsumerAdapter
	memoryCh chan *channels.Message
}

func NewMemmoryAdapter() *MemoryAdapter {
	return &MemoryAdapter{
		memoryCh: make(chan *channels.Message),
	}
}

func (ma *MemoryAdapter) Send(msg *channels.Message) error {
	ma.memoryCh <- msg
	return nil
}

func (ma *MemoryAdapter) Subscribe(client *channels.Client) error {
	ch := client.MessageChan()
	go func() {
		for {
			msg := <-ma.memoryCh
			ch <- msg
		}
	}()
	return nil
}
