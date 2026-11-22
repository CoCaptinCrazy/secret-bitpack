package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"secret-bitpack/internal/orchestration/kubernetes"
)

type assemblyRunner struct {
	args  []string
	input []byte
	err   error
}

func (r *assemblyRunner) Run(_ context.Context, _ string, args []string, input io.Reader, output, _ io.Writer) error {
	r.args = append([]string(nil), args...)
	if input != nil {
		r.input, _ = io.ReadAll(input)
	}
	_, _ = io.WriteString(output, "deployment/api configured\n")
	return r.err
}

func TestDeployBPackAppliesResourcesAndRecordsModules(t *testing.T) {
	root := t.TempDir()
	oldDataDirectory := dataDirectory
	dataDirectory = filepath.Join(root, "data")
	t.Cleanup(func() {
		dataDirectory = oldDataDirectory
	})
	runner := &assemblyRunner{}
	client := kubernetes.Client{Namespace: "test-space", Runner: runner}
	var output bytes.Buffer
	if err := deployBPack(context.Background(),
		filepath.Join("..", "..", "examples", "bitpack.bpack"),
		"test-space", false, client, &output); err != nil {
		t.Fatal(err)
	}
	if strings.Join(runner.args, " ") != "--namespace test-space apply -f -" {
		t.Fatalf("unexpected kubectl call: %#v", runner.args)
	}
	for _, expected := range []string{`"kind":"Deployment"`, `"kind":"Service"`, `"image":"social-api"`} {
		if !strings.Contains(string(runner.input), expected) {
			t.Errorf("manifest missing %q", expected)
		}
	}
	assemblyPath := filepath.Join(dataDirectory, "assemblies", "social-network", "assembly.json")
	state, err := os.ReadFile(assemblyPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(state), `"provider": "kubernetes"`) ||
		!strings.Contains(string(state), `"status": "active"`) {
		t.Fatalf("Kubernetes deployment state was not recorded: %s", state)
	}
	modules, err := moduleStore().List()
	if err != nil {
		t.Fatal(err)
	}
	if len(modules) != 6 || modules[0].Runtime != "kubernetes:test-space" {
		t.Fatalf("Kubernetes modules were not registered: %#v", modules)
	}
}

func TestDeployBPackFailureDoesNotRecordActiveAssembly(t *testing.T) {
	root := t.TempDir()
	oldDataDirectory := dataDirectory
	dataDirectory = filepath.Join(root, "data")
	t.Cleanup(func() {
		dataDirectory = oldDataDirectory
	})
	runner := &assemblyRunner{err: errors.New("kubectl unavailable")}
	client := kubernetes.Client{Namespace: "default", Runner: runner}
	err := deployBPack(context.Background(),
		filepath.Join("..", "..", "examples", "bitpack.bpack"),
		"default", false, client, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "kubectl unavailable") {
		t.Fatalf("expected Kubernetes apply failure, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(dataDirectory, "assemblies", "social-network", "assembly.json")); !os.IsNotExist(err) {
		t.Fatalf("failed deployment was recorded as active: %v", err)
	}
}

func TestBPackDryRunDoesNotCallKubernetes(t *testing.T) {
	runner := &assemblyRunner{}
	client := kubernetes.Client{Namespace: "default", Runner: runner}
	var output bytes.Buffer
	if err := deployBPack(context.Background(),
		filepath.Join("..", "..", "examples", "bitpack.bpack"),
		"default", true, client, &output); err != nil {
		t.Fatal(err)
	}
	if runner.args != nil {
		t.Fatalf("dry run contacted Kubernetes: %#v", runner.args)
	}
	if !strings.Contains(output.String(), `"kind": "Deployment"`) {
		t.Fatal("dry run did not print Kubernetes resources")
	}
}
