package alerts

import "sync"

type Manager struct {
	mu     sync.RWMutex
	alerts []Alert
}

func (m *Manager) Add(alert Alert) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alerts = append(m.alerts, alert)
}

func (m *Manager) List() []Alert {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]Alert(nil), m.alerts...)
}
