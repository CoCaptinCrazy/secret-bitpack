package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

var validClusterName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

type Store struct {
	Directory string
	Version   string
}

type Record struct {
	Kind       string   `json:"kind"`
	APIVersion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Status     Status   `json:"status"`
	Nodes      []any    `json:"nodes"`
	Assemblies []string `json:"assemblies"`
	Modules    []string `json:"modules"`
	Meshes     []string `json:"meshes"`
	Snapshots  []string `json:"snapshots"`
	Events     []Event  `json:"events"`
	Metrics    Metrics  `json:"metrics"`
}

type Metadata struct {
	Name           string `json:"name"`
	Created        int64  `json:"created"`
	BitpackVersion string `json:"bitpackVersion"`
}

type Status struct {
	Phase   string `json:"phase"`
	Healthy bool   `json:"healthy"`
}

type Event struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Metrics struct {
	Nodes      int `json:"nodes"`
	Assemblies int `json:"assemblies"`
	Modules    int `json:"modules"`
	Meshes     int `json:"meshes"`
	Snapshots  int `json:"snapshots"`
}

func NewStore(directory, version string) *Store {
	return &Store{Directory: directory, Version: version}
}

func (s *Store) Create(name string) (*Record, error) {
	if !validClusterName.MatchString(name) {
		return nil, fmt.Errorf("invalid cluster name %q; use letters, numbers, dashes, or underscores", name)
	}
	if err := os.MkdirAll(s.Directory, 0o700); err != nil {
		return nil, fmt.Errorf("create cluster directory: %w", err)
	}
	now := time.Now().Unix()
	record := &Record{
		Kind:       "Cluster",
		APIVersion: "v1",
		Metadata: Metadata{
			Name:           name,
			Created:        now,
			BitpackVersion: s.Version,
		},
		Status:     Status{Phase: "Pending", Healthy: false},
		Nodes:      []any{},
		Assemblies: []string{},
		Modules:    []string{},
		Meshes:     []string{},
		Snapshots:  []string{},
		Events: []Event{{
			Type:      "ClusterCreated",
			Message:   fmt.Sprintf("Cluster %s created", name),
			Timestamp: now,
		}},
		Metrics: Metrics{},
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode cluster record: %w", err)
	}
	path := s.path(name)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, os.ErrExist) {
		return nil, fmt.Errorf("cluster %q already exists", name)
	}
	if err != nil {
		return nil, fmt.Errorf("create cluster %q: %w", name, err)
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		os.Remove(path)
		return nil, fmt.Errorf("write cluster %q: %w", name, err)
	}
	if err := file.Close(); err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("close cluster %q record: %w", name, err)
	}
	return record, nil
}

func (s *Store) Get(name string) (*Record, error) {
	if !validClusterName.MatchString(name) {
		return nil, fmt.Errorf("invalid cluster name %q", name)
	}
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		return nil, fmt.Errorf("read cluster %q: %w", name, err)
	}
	var record Record
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("decode cluster %q: %w", name, err)
	}
	return &record, nil
}

func (s *Store) List() ([]Record, error) {
	entries, err := os.ReadDir(s.Directory)
	if errors.Is(err, os.ErrNotExist) {
		return []Record{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list clusters: %w", err)
	}
	records := make([]Record, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		name := entry.Name()[:len(entry.Name())-len(filepath.Ext(entry.Name()))]
		record, err := s.Get(name)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Metadata.Name < records[j].Metadata.Name
	})
	return records, nil
}

func (s *Store) Delete(name string) error {
	if !validClusterName.MatchString(name) {
		return fmt.Errorf("invalid cluster name %q", name)
	}
	if err := os.Remove(s.path(name)); err != nil {
		return fmt.Errorf("delete cluster %q: %w", name, err)
	}
	return nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.Directory, name+".json")
}
