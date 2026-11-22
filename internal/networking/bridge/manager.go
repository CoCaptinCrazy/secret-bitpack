package bridge

type Manager struct{ bridges map[string]Bridge }

func NewManager() *Manager { return &Manager{bridges: make(map[string]Bridge)} }

func (m *Manager) Add(bridge Bridge) { m.bridges[bridge.Name] = bridge }

func (m *Manager) Get(name string) (Bridge, bool) {
	bridge, ok := m.bridges[name]
	return bridge, ok
}
