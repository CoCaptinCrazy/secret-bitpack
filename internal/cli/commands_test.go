package cli

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/core"
)

func TestRegistryCommandsReachGoHandlers(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	commands := []struct {
		name  string
		args  []string
		input string
	}{
		{"help", nil, ""},
		{"version", nil, ""},
		{"forge", nil, ""},
		{"pack", nil, ""},
		{"unpack", nil, ""},
		{"verify", nil, ""},
		{"blueprints", nil, ""},
		{"inspect", nil, ""},
		{"export", nil, ""},
		{"import", nil, ""},
		{"ignite", nil, ""},
		{"modules", nil, ""},
		{"stop", nil, ""},
		{"restart", nil, ""},
		{"logs", nil, ""},
		{"snapshot", nil, ""},
		{"restore", nil, ""},
		{"assemble", nil, ""},
		{"assemblies", nil, ""},
		{"destroy", nil, ""},
		{"scale", nil, ""},
		{"rollout", nil, ""},
		{"rollback", nil, ""},
		{"publish", nil, ""},
		{"acquire", nil, ""},
		{"search", nil, ""},
		{"tags", nil, ""},
		{"login", nil, ""},
		{"logout", nil, ""},
		{"dev", []string{"inspect", "missing.bpdev"}, ""},
		{"shell", nil, "exit\n"},
		{"watch", []string{"one", "two"}, ""},
		{"batch", nil, ""},
		{"batch-status", nil, ""},
		{"mesh", nil, ""},
		{"meshes", nil, ""},
		{"bridges", nil, ""},
		{"route", nil, ""},
		{"health", nil, ""},
		{"policy", nil, ""},
		{"alert", nil, ""},
		{"metrics", nil, ""},
		{"profile", nil, ""},
		{"alerts", nil, ""},
		{"events", nil, ""},
		{"secrets", nil, ""},
		{"cluster", nil, ""},
		{"clusters", nil, ""},
		{"nodes", nil, ""},
		{"runtime", nil, ""},
		{"config", nil, ""},
		{"backup", nil, ""},
		{"pipeline", nil, ""},
		{"inspect-module", nil, ""},
		{"rm", nil, ""},
		{"plugin", nil, ""},
	}
	bitpack := core.New(io.Discard)
	for _, command := range commands {
		t.Run(command.name, func(t *testing.T) {
			err := execute(bitpack, append([]string{command.name}, command.args...),
				strings.NewReader(command.input), io.Discard)
			if err != nil && strings.Contains(err.Error(), "unknown command") {
				t.Fatalf("registered command has no handler: %v", err)
			}
		})
	}
}

func TestForgePreservesPolyglotBlueprintFields(t *testing.T) {
	root := t.TempDir()
	oldDataDirectory := dataDirectory
	dataDirectory = filepath.Join(root, "data")
	t.Cleanup(func() {
		dataDirectory = oldDataDirectory
	})
	bitfile := filepath.Join(root, "Bitfile")
	content := `blueprint "polyglot-app" v1
using application
bring "./src"
env <|
    APP_MODE = production
|>
secret <|
    signing-key
|>
assemble <|
    cargo build --release
|>
verify <|
    go test ./...
|>
ignite <|
    node server.js
|>
`
	if err := os.WriteFile(bitfile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := bitpackForge(io.Discard, []string{bitfile}); err != nil {
		t.Fatal(err)
	}

	document, err := readBlueprint(filepath.Join(dataDirectory, "blueprints", "polyglot-app.bp"))
	if err != nil {
		t.Fatal(err)
	}
	runtime, ok := document["runtime"].(map[string]any)
	if !ok || runtime["using"] != "application" {
		t.Fatalf("runtime declaration was not retained: %#v", document["runtime"])
	}
	environment, ok := runtime["environment"].(map[string]any)
	if !ok || environment["APP_MODE"] != "production" {
		t.Fatalf("environment was not retained: %#v", runtime["environment"])
	}
	execution, ok := document["execution"].(map[string]any)
	if !ok {
		t.Fatalf("execution instructions were not retained: %#v", document["execution"])
	}
	assemble, _ := execution["assemble"].([]any)
	verify, _ := execution["verify"].([]any)
	ignite, _ := execution["ignite"].([]any)
	if len(assemble) != 1 || assemble[0] != "cargo build --release" ||
		len(verify) != 1 || verify[0] != "go test ./..." ||
		len(ignite) != 1 || ignite[0] != "node server.js" {
		t.Fatalf("polyglot commands were not retained: %#v", document["execution"])
	}
	security, ok := document["security"].(map[string]any)
	if !ok {
		t.Fatalf("security metadata was not retained: %#v", document["security"])
	}
	secrets, _ := security["secrets"].([]any)
	if len(secrets) != 1 || secrets[0] != "signing-key" {
		t.Fatalf("secret names were not retained: %#v", document["security"])
	}
	source, ok := document["source"].(map[string]any)
	if !ok || source["context"] != "./src" {
		t.Fatalf("project context was not retained: %#v", document["source"])
	}

	encoded, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(encoded) {
		t.Fatal("forged blueprint is not valid JSON")
	}
}

func TestBlueprintPhasesRunGenericShellCommands(t *testing.T) {
	root := t.TempDir()
	oldDataDirectory := dataDirectory
	dataDirectory = filepath.Join(root, "data")
	t.Cleanup(func() {
		dataDirectory = oldDataDirectory
	})
	bitfile := filepath.Join(root, "Bitfile")
	content := `blueprint "polyglot-app" v1
using application
bring "./"
assemble <|
    echo assembled
|>
verify <|
    echo verified
|>
ignite <|
    echo launched
|>
`
	if err := os.WriteFile(bitfile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := bitpackForge(io.Discard, []string{bitfile}); err != nil {
		t.Fatal(err)
	}

	var output strings.Builder
	if err := bitpackAssemble(&output, []string{"polyglot-app"}); err != nil {
		t.Fatal(err)
	}
	if err := bitpackVerify(&output, []string{"polyglot-app"}); err != nil {
		t.Fatal(err)
	}
	if err := bitpackIgnite(&output, []string{"polyglot-app"}); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"assembled", "verified", "Background process started"} {
		if !strings.Contains(output.String(), expected) {
			t.Errorf("expected command output %q, got %q", expected, output.String())
		}
	}
	lastIndex := -1
	for _, expected := range []string{"assembled", "verified", "Background process started"} {
		index := strings.Index(output.String(), expected)
		if index <= lastIndex {
			t.Errorf("commands did not execute in order: %q", output.String())
			break
		}
		lastIndex = index
	}
	module, err := findLocalModule(moduleStore(), "polyglot-app")
	if err != nil || module == nil {
		t.Fatalf("find started module: module=%#v err=%v", module, err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		log, err := os.ReadFile(filepath.Join(dataDirectory, "logs", "modules", module.ID+".log"))
		if err == nil && strings.Contains(string(log), "launched") {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("background ignite command did not write its output")
}

func TestBlueprintVerifyReturnsCommandFailures(t *testing.T) {
	root := t.TempDir()
	bitfile := filepath.Join(root, "Bitfile")
	content := `blueprint "failing-app" v1
bring "./"
verify <|
    exit 7
|>
`
	if err := os.WriteFile(bitfile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := bitpackVerify(io.Discard, []string{bitfile}); err == nil {
		t.Fatal("expected verification command failure to be returned")
	}
}
