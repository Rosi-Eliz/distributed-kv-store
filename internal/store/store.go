package store

import (
	"encoding/json"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// SerializeData serializes the current data map to a JSON byte slice.
func (s *Store) SerializeData() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.data)
}

// ReplaceData replaces the current data map with the provided data map.
func (s *Store) ReplaceData(newdata map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = newdata
}
