package manager

import (
	"sync"
	"time"
)

type Supervisor struct {
	mu       sync.RWMutex
	lastSeen map[string]time.Time
}

func NewSupervisor() *Supervisor { return &Supervisor{lastSeen: make(map[string]time.Time)} }

func (s *Supervisor) Heartbeat(id string, at time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastSeen[id] = at
}

func (s *Supervisor) LastSeen(id string) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.lastSeen[id]
	return value, ok
}
