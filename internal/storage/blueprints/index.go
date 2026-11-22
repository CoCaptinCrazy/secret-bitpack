package blueprints

import "sync"

type Index struct {
	mu      sync.RWMutex
	entries map[string]string
}

func NewIndex() *Index { return &Index{entries: make(map[string]string)} }

func (i *Index) Set(name, path string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.entries[name] = path
}

func (i *Index) Get(name string) (string, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	path, ok := i.entries[name]
	return path, ok
}
