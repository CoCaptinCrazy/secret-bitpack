package manager

import (
	"sync"

	"secret-bitpack/internal/runtime/modules"
)

type Manager struct {
	mu      sync.RWMutex
	modules map[string]*modules.Module
}

func New() *Manager {
	return &Manager{modules: make(map[string]*modules.Module)}
}

func (m *Manager) Start(module *modules.Module) {
	m.mu.Lock()
	defer m.mu.Unlock()
	module.Status = "running"
	m.modules[module.ID] = module
}

func (m *Manager) Stop(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	module, ok := m.modules[id]
	if ok {
		module.Status = "stopped"
	}
	return ok
}

func (m *Manager) Status(id string) (*modules.Module, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	module, ok := m.modules[id]
	return module, ok
}

func (m *Manager) List() []*modules.Module {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*modules.Module, 0, len(m.modules))
	for _, module := range m.modules {
		result = append(result, module)
	}
	return result
}
