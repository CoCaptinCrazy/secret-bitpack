package storage

import (
	"os"
	"path/filepath"
)

type Manager struct {
	Root string
}

func New(root string) *Manager {
	return &Manager{Root: root}
}

func (m *Manager) Initialize() error {
	for _, name := range []string{"blueprints", "modules", "snapshots", "secrets"} {
		if err := os.MkdirAll(filepath.Join(m.Root, name), 0o700); err != nil {
			return err
		}
	}
	return nil
}
