package certificates

import "sync"

type Store struct {
	mu           sync.RWMutex
	certificates map[string]Certificate
}

func NewStore() *Store { return &Store{certificates: make(map[string]Certificate)} }

func (s *Store) Put(name string, certificate Certificate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.certificates[name] = certificate
}

func (s *Store) Get(name string) (Certificate, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	certificate, ok := s.certificates[name]
	return certificate, ok
}
