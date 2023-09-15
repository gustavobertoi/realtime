package channels

type ClientMemoryStore struct {
	ClientStore
	clientMap map[string]*Client
}

func NewClientMemoryStore() *ClientMemoryStore {
	return &ClientMemoryStore{
		clientMap: make(map[string]*Client),
	}
}

func (s *ClientMemoryStore) Count() int {
	return len(s.clientMap)
}

func (s *ClientMemoryStore) All() []*Client {
	var clients []*Client
	for _, client := range s.clientMap {
		clients = append(clients, client)
	}
	return clients
}

func (s *ClientMemoryStore) Get(id string) (*Client, error) {
	if !s.Has(id) {
		return nil, errMessageDoesNotExist
	}
	return s.clientMap[id], nil
}

func (s *ClientMemoryStore) Has(id string) bool {
	return s.clientMap[id] != nil
}

func (s *ClientMemoryStore) Put(c *Client) error {
	s.clientMap[c.ID] = c
	return nil
}

func (s *ClientMemoryStore) Delete(c *Client) error {
	if s.Has(c.ID) {
		delete(s.clientMap, c.ID)
	}
	return nil
}
