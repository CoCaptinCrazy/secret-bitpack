package modules

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

var validModuleName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`)

type Store struct {
	Directory string
}

func NewStore(directory string) *Store {
	return &Store{Directory: directory}
}

func (s *Store) Create(name, blueprint, runtime string) (*Module, error) {
	if !validModuleName.MatchString(name) {
		return nil, fmt.Errorf("invalid module name %q", name)
	}
	if err := os.MkdirAll(s.Directory, 0o700); err != nil {
		return nil, fmt.Errorf("create module directory: %w", err)
	}
	random := make([]byte, 4)
	if _, err := rand.Read(random); err != nil {
		return nil, fmt.Errorf("generate module ID: %w", err)
	}
	module := &Module{
		ID: "mod-" + hex.EncodeToString(random), Name: name, Status: "running",
		Blueprint: blueprint, Runtime: runtime, Started: time.Now().Unix(),
	}
	if err := s.write(module); err != nil {
		return nil, err
	}
	return module, nil
}

func (s *Store) List() ([]Module, error) {
	entries, err := os.ReadDir(s.Directory)
	if os.IsNotExist(err) {
		return []Module{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list modules: %w", err)
	}
	var result []Module
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		var module Module
		data, err := os.ReadFile(filepath.Join(s.Directory, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read module %s: %w", entry.Name(), err)
		}
		if err := json.Unmarshal(data, &module); err != nil {
			return nil, fmt.Errorf("decode module %s: %w", entry.Name(), err)
		}
		result = append(result, module)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *Store) Find(identifier string) (*Module, error) {
	if validModuleName.MatchString(identifier) && len(identifier) > 4 && identifier[:4] == "mod-" {
		module, err := s.read(identifier)
		if err == nil {
			return module, nil
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	modules, err := s.List()
	if err != nil {
		return nil, err
	}
	for i := range modules {
		if modules[i].ID == identifier || modules[i].Name == identifier {
			return &modules[i], nil
		}
	}
	return nil, fmt.Errorf("module %q not found", identifier)
}

func (s *Store) Update(module *Module) error {
	return s.write(module)
}

func (s *Store) Delete(module *Module) error {
	if err := os.Remove(s.path(module.ID)); err != nil {
		return fmt.Errorf("delete module %q: %w", module.ID, err)
	}
	return nil
}

func (s *Store) Path(module *Module) string {
	return s.path(module.ID)
}

func (s *Store) write(module *Module) error {
	data, err := json.MarshalIndent(module, "", "  ")
	if err != nil {
		return fmt.Errorf("encode module: %w", err)
	}
	if err := os.MkdirAll(s.Directory, 0o700); err != nil {
		return fmt.Errorf("create module directory: %w", err)
	}
	if err := os.WriteFile(s.path(module.ID), data, 0o600); err != nil {
		return fmt.Errorf("write module %q: %w", module.ID, err)
	}
	return nil
}

func (s *Store) read(id string) (*Module, error) {
	data, err := os.ReadFile(s.path(id))
	if err != nil {
		return nil, err
	}
	var module Module
	if err := json.Unmarshal(data, &module); err != nil {
		return nil, fmt.Errorf("decode module %q: %w", id, err)
	}
	return &module, nil
}

func (s *Store) path(id string) string {
	return filepath.Join(s.Directory, id+".json")
}
