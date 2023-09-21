package channels

import (
	"github.com/open-source-cloud/realtime/pkg/store"
)

type ClientStore struct {
	store *store.MemoryStore
}

func NewClientStore() *ClientStore {
	return &ClientStore{
		store: store.NewMemoryStore(),
	}
}

func (cs *ClientStore) Count() int {
	return cs.store.Count()
}

func (cs *ClientStore) All() []*Client {
	var clients []*Client
	cs.store.ForEach(func(_ string, value interface{}) {
		clients = append(clients, value.(*Client))
	})
	return clients
}

func (cs *ClientStore) Get(id string) (*Client, error) {
	if !cs.Has(id) {
		return nil, errMessageDoesNotExist
	}
	client, err := cs.store.Get(id)
	if err != nil {
		return nil, err
	}
	return client.(*Client), nil
}

func (cs *ClientStore) Has(id string) bool {
	return cs.store.Has(id)
}

func (cs *ClientStore) Put(c *Client) {
	cs.store.Set(c.ID, c)
}

func (cs *ClientStore) Delete(c *Client) {
	cs.store.Delete(c.ID)
}
