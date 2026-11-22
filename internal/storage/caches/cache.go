package caches

import "sync"

type Cache struct {
	mu     sync.RWMutex
	values map[string][]byte
}

func New() *Cache { return &Cache{values: make(map[string][]byte)} }

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = append([]byte(nil), value...)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.values[key]
	return append([]byte(nil), value...), ok
}
