package state

import (
	"os"
	"sync"
)

type Manager struct {
	mu   sync.RWMutex
	data map[string]string
}

func New() *Manager {
	return &Manager{data: make(map[string]string)}
}

func (m *Manager) Load() error {
	return os.MkdirAll(".bitpack", 0o755)
}

func (m *Manager) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *Manager) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.data[key]
	return value, ok
}

func (m *Manager) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}
