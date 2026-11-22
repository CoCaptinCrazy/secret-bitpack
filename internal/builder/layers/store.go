package layers

import "sync"

type Store struct {
	mu     sync.RWMutex
	layers map[[32]byte][]byte
}

func NewStore() *Store { return &Store{layers: make(map[[32]byte][]byte)} }

func (s *Store) Put(content []byte) Layer {
	layer := New(content)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.layers[layer.Digest] = append([]byte(nil), content...)
	return layer
}

func (s *Store) Get(digest [32]byte) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	content, ok := s.layers[digest]
	return append([]byte(nil), content...), ok
}
