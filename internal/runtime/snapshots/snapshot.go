package snapshots

import "time"

type Snapshot struct {
	ID        string
	ModuleID  string
	CreatedAt time.Time
}

type Store struct {
	Snapshots []Snapshot
}

func (s *Store) Create(moduleID string) Snapshot {
	snapshot := Snapshot{
		ID:        moduleID + "-" + time.Now().UTC().Format("20060102T150405Z"),
		ModuleID:  moduleID,
		CreatedAt: time.Now().UTC(),
	}
	s.Snapshots = append(s.Snapshots, snapshot)
	return snapshot
}
