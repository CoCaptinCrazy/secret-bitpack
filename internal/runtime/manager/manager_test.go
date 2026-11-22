package manager

import (
	"testing"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/runtime/modules"
)

func TestManagerLifecycle(t *testing.T) {
	manager := New()
	module := &modules.Module{ID: "api", Name: "social-api", Version: "1.4"}
	manager.Start(module)

	status, ok := manager.Status("api")
	if !ok || status.Status != "running" {
		t.Fatal("module did not start")
	}
	if !manager.Stop("api") {
		t.Fatal("module did not stop")
	}
	status, _ = manager.Status("api")
	if status.Status != "stopped" {
		t.Fatal("module did not stop")
	}
}
