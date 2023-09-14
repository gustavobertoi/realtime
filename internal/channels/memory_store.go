package channels

type MessageMemoryStore struct {
	MessageStore
	msgMap map[string]*Message
}

func NewMessageMemoryStore() *MessageMemoryStore {
	return &MessageMemoryStore{
		msgMap: make(map[string]*Message),
	}
}

func (store *MessageMemoryStore) Get(id string) (*Message, error) {
	if !store.Has(id) {
		return nil, errMessageDoesNotExist
	}
	return store.msgMap[id], nil
}

func (store *MessageMemoryStore) Has(id string) bool {
	return store.msgMap[id] != nil
}

func (store *MessageMemoryStore) Put(msg *Message) error {
	store.msgMap[msg.ID] = msg
	return nil
}

func (store *MessageMemoryStore) Delete(msg *Message) error {
	if store.Has(msg.ID) {
		delete(store.msgMap, msg.ID)
	}
	return nil
}
