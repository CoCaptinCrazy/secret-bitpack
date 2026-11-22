package secrets

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Manager struct {
	mu      sync.RWMutex
	secrets map[string]string
}

func New() *Manager {
	return &Manager{secrets: make(map[string]string)}
}

func (m *Manager) Set(name, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.secrets[name] = value
}

func (m *Manager) Get(name string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.secrets[name]
	return value, ok
}

func (m *Manager) Delete(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.secrets, name)
}

func EnvironmentKey(name string) string {
	var key strings.Builder
	for _, char := range strings.ToUpper(name) {
		if char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' {
			key.WriteRune(char)
		} else {
			key.WriteByte('_')
		}
	}
	return "BITPACK_SECRET_" + key.String()
}

func ResolveFromEnvironment(names []string) (map[string]string, error) {
	values := make(map[string]string, len(names))
	usedKeys := make(map[string]string, len(names))
	for _, name := range names {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("secret name cannot be empty")
		}
		environmentKey := EnvironmentKey(name)
		if previous, exists := usedKeys[environmentKey]; exists && previous != name {
			return nil, fmt.Errorf("secret names %q and %q map to the same environment variable", previous, name)
		}
		usedKeys[environmentKey] = name
		value, ok := os.LookupEnv(environmentKey)
		if !ok {
			return nil, fmt.Errorf("secret %q is missing; set %s", name, environmentKey)
		}
		values[name] = value
	}
	return values, nil
}
