package metrics

import (
	"sync"
	"sync/atomic"
)

type Registry struct {
	mu     sync.RWMutex
	values map[string]*atomic.Int64
}

func New() *Registry {
	return &Registry{values: make(map[string]*atomic.Int64)}
}

func (r *Registry) Add(name string, value int64) {
	r.mu.Lock()
	counter, ok := r.values[name]
	if !ok {
		counter = &atomic.Int64{}
		r.values[name] = counter
	}
	r.mu.Unlock()
	counter.Add(value)
}

func (r *Registry) Get(name string) int64 {
	r.mu.RLock()
	counter, ok := r.values[name]
	r.mu.RUnlock()
	if !ok {
		return 0
	}
	return counter.Load()
}
