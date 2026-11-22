package modules

import "sync"

type Database struct {
	mu      sync.RWMutex
	records map[string]Metadata
}

func NewDatabase() *Database { return &Database{records: make(map[string]Metadata)} }

func (d *Database) Put(metadata Metadata) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.records[metadata.ID] = metadata
}

func (d *Database) Get(id string) (Metadata, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	metadata, ok := d.records[id]
	return metadata, ok
}
