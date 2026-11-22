package dns

import "strings"

type Resolver struct{ records map[string][]Record }

func NewResolver() *Resolver { return &Resolver{records: make(map[string][]Record)} }

func (r *Resolver) Add(record Record) {
	key := strings.ToLower(strings.TrimSuffix(record.Name, "."))
	r.records[key] = append(r.records[key], record)
}

func (r *Resolver) Lookup(name string) []Record {
	key := strings.ToLower(strings.TrimSuffix(name, "."))
	return append([]Record(nil), r.records[key]...)
}
