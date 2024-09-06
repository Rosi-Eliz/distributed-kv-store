package store

import (
	"encoding/json"
	"sync"
	"time"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
	ttl  map[string]time.Time
}

func NewStore() *Store {
	s := &Store{
		data: make(map[string]string),
		ttl:  make(map[string]time.Time),
	}
	go s.backgroundCleaner() // Start the background cleaner
	return s
}

func (s *Store) backgroundCleaner() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, exp := range s.ttl {
			if now.After(exp) {
				delete(s.data, key)
				delete(s.ttl, key)
			}
		}
		s.mu.Unlock()
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	// If a TTL is provided, set the expiration time
	if ttl > 0 {
		s.ttl[key] = time.Now().Add(ttl)
	} else {
		delete(s.ttl, key) // No TTL means no expiration
	}
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	delete(s.ttl, key)
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
