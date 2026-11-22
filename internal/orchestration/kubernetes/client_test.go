package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

type fakeRunner struct {
	args   []string
	calls  [][]string
	input  []byte
	err    error
	stderr string
}

func (r *fakeRunner) Run(_ context.Context, _ string, args []string, input io.Reader, output, errorOutput io.Writer) error {
	r.args = append([]string(nil), args...)
	r.calls = append(r.calls, append([]string(nil), args...))
	if input != nil {
		r.input, _ = io.ReadAll(input)
	}
	if r.stderr != "" {
		_, _ = io.WriteString(errorOutput, r.stderr)
	}
	_, _ = io.WriteString(output, "deployment/api configured\n")
	return r.err
}

func TestApplyUsesKubeconfigAndDoesNotLeakSecrets(t *testing.T) {
	runner := &fakeRunner{err: errors.New("kubectl exited with status 1"), stderr: "rejected test-secret-value"}
	client := Client{Namespace: "test-space", Runner: runner}
	objects := []map[string]any{{
		"apiVersion": "v1", "kind": "Secret",
		"stringData": map[string]string{"password": "test-secret-value"},
	}}
	var output bytes.Buffer
	err := client.Apply(context.Background(), objects,
		map[string]string{"password": "test-secret-value"}, &output)
	if err == nil {
		t.Fatal("expected kubectl failure")
	}
	if strings.Contains(err.Error(), "test-secret-value") || strings.Contains(output.String(), "test-secret-value") {
		t.Fatal("secret value was exposed in diagnostics")
	}
	if !strings.Contains(err.Error(), "[REDACTED]") {
		t.Fatalf("kubectl diagnostic was not retained with redaction: %v", err)
	}
	if strings.Join(runner.args, " ") != "--namespace test-space apply -f -" {
		t.Fatalf("unexpected kubectl invocation: %#v", runner.args)
	}
	if !bytes.Contains(runner.input, []byte("test-secret-value")) {
		t.Fatal("secret was not sent to kubectl")
	}
}

func TestScaleUsesArgumentsWithoutShell(t *testing.T) {
	runner := &fakeRunner{}
	client := Client{Namespace: "default", Runner: runner}
	if err := client.Scale(context.Background(), "social-api", 0, io.Discard); err != nil {
		t.Fatal(err)
	}
	expected := []string{"--namespace", "default", "scale", "deployment/social-api", "--replicas=0"}
	if strings.Join(runner.args, " ") != strings.Join(expected, " ") {
		t.Fatalf("unexpected kubectl invocation: %#v", runner.args)
	}
}

func TestStopRemovesAutoscalerBeforeScalingDown(t *testing.T) {
	runner := &fakeRunner{}
	client := Client{Namespace: "default", Runner: runner}
	if err := client.Stop(context.Background(), "social-api", io.Discard); err != nil {
		t.Fatal(err)
	}
	if len(runner.calls) != 2 {
		t.Fatalf("expected autoscaler deletion and scale operations, got %#v", runner.calls)
	}
	if got, want := strings.Join(runner.calls[0], " "),
		"--namespace default delete hpa/social-api --ignore-not-found=true"; got != want {
		t.Fatalf("unexpected autoscaler operation %q, want %q", got, want)
	}
	if got, want := strings.Join(runner.calls[1], " "),
		"--namespace default scale deployment/social-api --replicas=0"; got != want {
		t.Fatalf("unexpected scale operation %q, want %q", got, want)
	}
}
