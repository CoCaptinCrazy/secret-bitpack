package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestInfoWritesMessage(t *testing.T) {
	var output bytes.Buffer
	New(&output).Info("started %s", "bitpack")
	if !strings.Contains(output.String(), "[info] started bitpack") {
		t.Fatalf("unexpected log output %q", output.String())
	}
}
