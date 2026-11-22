package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseExampleBPack(t *testing.T) {
	path := filepath.Join("..", "..", "..", "examples", "bitpack.bpack")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	assembly, err := ParseBPack(file)
	if err != nil {
		t.Fatal(err)
	}
	if assembly.Name != "social-network" || len(assembly.Meshes) != 2 ||
		len(assembly.Modules) != 6 || len(assembly.Rules) != 5 || len(assembly.Autoscale) != 2 {
		t.Fatalf("unexpected assembly: %#v", assembly)
	}
	api := findModule(assembly.Modules, "api")
	if api == nil || api.Blueprint != "social-api" || api.Replicas != 4 ||
		len(api.Volumes) != 2 || len(api.Depends) != 2 {
		t.Fatalf("API module was not parsed: %#v", api)
	}
	if api.Volumes[0].Size != "1Gi" {
		t.Fatalf("volume default size was not set: %#v", api.Volumes[0])
	}
}

func TestParseBPackRejectsUnknownDependency(t *testing.T) {
	_, err := ParseBPack(strings.NewReader(`assembly app {
    module api {
        blueprint example/api:1
        depends missing
    }
}`))
	if err == nil || !strings.Contains(err.Error(), "unknown module") {
		t.Fatalf("expected unknown dependency error, got %v", err)
	}
}

func TestParseBPackRejectsUnclosedBlock(t *testing.T) {
	if _, err := ParseBPack(strings.NewReader("assembly app {\n")); err == nil {
		t.Fatal("expected an unclosed block error")
	}
}

func findModule(modules []Module, name string) *Module {
	for index := range modules {
		if modules[index].Name == name {
			return &modules[index]
		}
	}
	return nil
}
