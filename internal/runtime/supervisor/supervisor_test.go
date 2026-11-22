package supervisor

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestManagerStartsStopsAndRestartsProcess(t *testing.T) {
	root := t.TempDir()
	manager := New(filepath.Join(root, "state"), filepath.Join(root, "logs"))
	spec := Spec{
		Name: "test-process", Executable: os.Args[0],
		Arguments:  []string{"-test.run=^TestSupervisorHelperProcess$"},
		WorkingDir: root, Environment: map[string]string{"BITPACK_SUPERVISOR_HELPER": "1"},
	}

	record, err := manager.Start("test-process", spec)
	if err != nil {
		t.Fatal(err)
	}
	if running, err := manager.IsRunning(record.ID); err != nil || !running {
		t.Fatalf("new process not running: running=%t err=%v", running, err)
	}
	if err := manager.Stop(record.ID); err != nil {
		t.Fatal(err)
	}
	if running, err := manager.IsRunning(record.ID); err != nil || running {
		t.Fatalf("stopped process still running: running=%t err=%v", running, err)
	}

	restarted, err := manager.Restart(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if running, err := manager.IsRunning(restarted.ID); err != nil || !running {
		t.Fatalf("restarted process not running: running=%t err=%v", running, err)
	}
	if err := manager.Stop(restarted.ID); err != nil {
		t.Fatal(err)
	}
}

func TestSupervisorHelperProcess(t *testing.T) {
	if os.Getenv("BITPACK_SUPERVISOR_HELPER") != "1" {
		return
	}
	for {
		time.Sleep(time.Second)
	}
}

func TestManagerDoesNotPersistSecretValues(t *testing.T) {
	root := t.TempDir()
	manager := New(filepath.Join(root, "state"), filepath.Join(root, "logs"))
	t.Setenv("BITPACK_SECRET_API_TOKEN", "must-not-be-saved")
	spec := Spec{
		Name: "secret-process", Executable: os.Args[0],
		Arguments:  []string{"-test.run=^TestSupervisorHelperProcess$"},
		WorkingDir: root, Environment: map[string]string{"BITPACK_SUPERVISOR_HELPER": "1"},
		Secrets: []string{"api-token"},
	}
	record, err := manager.Start("secret-process", spec)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := manager.Stop(record.ID); err != nil {
			t.Errorf("stop helper process: %v", err)
		}
	})
	state, err := os.ReadFile(filepath.Join(manager.Directory, record.ID+".json"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(state, []byte("must-not-be-saved")) {
		t.Fatal("secret value was persisted in process state")
	}
}
