package pubsub

import (
	"time"

	"github.com/open-source-cloud/realtime/channels"
)

type MemoryAdapter struct {
	channels.ProducerAdapter
	channels.ConsumerAdapter
	msgCh   chan *channels.Message
	clients []*channels.Client
}

func NewMemmoryAdapter() *MemoryAdapter {
	ma := &MemoryAdapter{
		msgCh: make(chan *channels.Message),
	}
	go ma.sendMessageToAllClientsHandler()
	return ma
}

func (ma *MemoryAdapter) Send(msg *channels.Message) error {
	ma.msgCh <- msg
	return nil
}

func (ma *MemoryAdapter) Subscribe(client *channels.Client) error {
	ma.clients = append(ma.clients, client)
	return nil
}

func (ma *MemoryAdapter) sendMessageToAllClientsHandler() {
	for {
		msg := <-ma.msgCh
		for _, client := range ma.clients {
			ch := client.MessageChan()
			ch <- msg
		}
		time.Sleep(2 * time.Second)
	}
}
