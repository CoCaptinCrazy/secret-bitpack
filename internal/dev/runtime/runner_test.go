package runtime

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"secret-bitpack/internal/dev/parser"
	"secret-bitpack/internal/runtime/supervisor"
)

func TestRunSupervisesServicesAndRestartsWatchedCommands(t *testing.T) {
	root := t.TempDir()
	watched := filepath.Join(root, "src", "main.txt")
	if err := os.MkdirAll(filepath.Dir(watched), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(watched, []byte("first"), 0o600); err != nil {
		t.Fatal(err)
	}
	command := helperCommand()
	environment := &parser.Environment{
		Name: "test-dev", Workspace: root,
		Services: []parser.Service{{Name: "helper", Command: command}},
		Watch:    []string{"src/"},
		Shell:    []string{command},
		Variables: map[string]string{
			"BITPACK_DEV_HELPER": "1",
		},
	}
	manager := supervisor.New(filepath.Join(root, "state"), filepath.Join(root, "logs"))
	ctx, cancel := context.WithCancel(context.Background())
	var output lockedBuffer
	result := make(chan error, 1)
	go func() {
		result <- Run(ctx, environment, root, manager, strings.NewReader(""), &output, io.Discard)
	}()

	if !waitForOutput(&output, "Development process started", 4*time.Second) {
		cancel()
		t.Fatal("development command did not start")
	}
	if err := os.WriteFile(watched, []byte("changed size"), 0o600); err != nil {
		cancel()
		t.Fatal(err)
	}
	if !waitForOutput(&output, "[RESTART]", 4*time.Second) {
		cancel()
		t.Fatal("watched file change did not restart the development command")
	}
	cancel()
	select {
	case err := <-result:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("development runner did not stop after cancellation")
	}
	for _, kindName := range []struct {
		kind string
		name string
	}{{"service", "helper"}, {"app", "shell"}} {
		id := processID(environment.Name, kindName.kind, 0, kindName.name)
		running, err := manager.IsRunning(id)
		if err != nil || running {
			t.Errorf("process %s was not stopped: running=%t err=%v", id, running, err)
		}
	}
}

func TestRunRequiresServiceCommand(t *testing.T) {
	environment := &parser.Environment{
		Name: "test", Services: []parser.Service{{Name: "database"}},
	}
	err := Run(context.Background(), environment, t.TempDir(),
		supervisor.New(filepath.Join(t.TempDir(), "state"), filepath.Join(t.TempDir(), "logs")),
		strings.NewReader(""), io.Discard, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "requires a host command") {
		t.Fatalf("expected a missing service command error, got %v", err)
	}
}

func TestDevHelperProcess(t *testing.T) {
	if os.Getenv("BITPACK_DEV_HELPER") != "1" {
		return
	}
	for {
		time.Sleep(time.Second)
	}
}

func helperCommand() string {
	executable := os.Args[0]
	if runtime.GOOS == "windows" {
		executable = `"` + strings.ReplaceAll(executable, `"`, `\"`) + `"`
	} else {
		executable = "'" + strings.ReplaceAll(executable, "'", "'\\''") + "'"
	}
	return fmt.Sprintf("%s -test.run=^TestDevHelperProcess$", executable)
}

type lockedBuffer struct {
	mu sync.Mutex
	b  strings.Builder
}

func (b *lockedBuffer) Write(value []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Write(value)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.String()
}

func waitForOutput(output *lockedBuffer, value string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if strings.Contains(output.String(), value) {
			return true
		}
		time.Sleep(25 * time.Millisecond)
	}
	return strings.Contains(output.String(), value)
}
