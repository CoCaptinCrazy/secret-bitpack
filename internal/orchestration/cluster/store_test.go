package cluster

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreCreateGetListAndDelete(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "clusters"), "test")
	created, err := store.Create("social-api")
	if err != nil {
		t.Fatal(err)
	}
	if created.Metadata.Name != "social-api" || created.Status.Phase != "Pending" || created.Status.Healthy {
		t.Fatalf("unexpected cluster record: %#v", created)
	}

	loaded, err := store.Get("social-api")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Metadata.BitpackVersion != "test" {
		t.Fatalf("unexpected version: %s", loaded.Metadata.BitpackVersion)
	}

	listed, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 || listed[0].Metadata.Name != "social-api" {
		t.Fatalf("unexpected cluster list: %#v", listed)
	}

	if err := store.Delete("social-api"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(store.Directory, "social-api.json")); !os.IsNotExist(err) {
		t.Fatalf("cluster record still exists or stat failed: %v", err)
	}
}

func TestStoreRejectsExistingAndUnsafeNames(t *testing.T) {
	store := NewStore(t.TempDir(), "test")
	if _, err := store.Create("social-api"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Create("social-api"); err == nil {
		t.Fatal("expected duplicate cluster creation to fail")
	}
	if _, err := store.Create("../outside"); err == nil {
		t.Fatal("expected unsafe cluster name to fail")
	}
}
