package modules

import (
	"path/filepath"
	"testing"
)

func TestStoreLifecycle(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "modules"))
	created, err := store.Create("social-api", "social-api", "node20")
	if err != nil {
		t.Fatal(err)
	}
	found, err := store.Find(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.Name != "social-api" || found.Status != "running" {
		t.Fatalf("unexpected module: %#v", found)
	}
	found.Status = "stopped"
	if err := store.Update(found); err != nil {
		t.Fatal(err)
	}
	listed, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 || listed[0].Status != "stopped" {
		t.Fatalf("module update was not persisted: %#v", listed)
	}
	if err := store.Delete(found); err != nil {
		t.Fatal(err)
	}
}
