package store

import "errors"

var (
	errItemDoesNotExists = errors.New("item does not exists in store")
)

type MemoryStore struct {
	data map[string]interface{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]interface{}),
	}
}

func (m *MemoryStore) Count() int {
	return len(m.data)
}

func (m *MemoryStore) ForEach(callback func(key string, value interface{})) {
	for key, value := range m.data {
		callback(key, value)
	}
}

func (m *MemoryStore) Get(key string) (interface{}, error) {
	if !m.Has(key) {
		return nil, errItemDoesNotExists
	}
	return m.data[key], nil
}

func (m *MemoryStore) Has(key string) bool {
	return m.data[key] != nil
}

func (m *MemoryStore) Set(key string, data interface{}) error {
	m.data[key] = m
	return nil
}

func (m *MemoryStore) Delete(key string) error {
	if m.Has(key) {
		delete(m.data, key)
	}
	return nil
}
