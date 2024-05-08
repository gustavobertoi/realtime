package pubsub

import (
	"time"

	"github.com/gustavobertoi/realtime/internal/channels"
)

type MemoryAdapter struct {
	msgCh   chan *channels.Message
	clients []*channels.Client
}

func NewMemoryAdapter() (*PubSub, error) {
	ma := &MemoryAdapter{
		msgCh: make(chan *channels.Message),
	}
	go ma.sendMessageToAllClientsHandler()
	return &PubSub{
		Driver:   MemoryDriver,
		Consumer: ma,
		Producer: ma,
	}, nil
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
