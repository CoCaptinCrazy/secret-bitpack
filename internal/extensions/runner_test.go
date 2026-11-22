package extensions

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunRejectsInvalidExtensionName(t *testing.T) {
	err := Run(bytes.NewReader(nil), &bytes.Buffer{}, &bytes.Buffer{}, t.TempDir(), "..\\outside")
	if err == nil {
		t.Fatal("expected path-like extension name to be rejected")
	}
}

func TestRunReportsMissingExtension(t *testing.T) {
	directory := t.TempDir()
	err := Run(bytes.NewReader(nil), &bytes.Buffer{}, &bytes.Buffer{}, directory, "missing")
	if err == nil {
		t.Fatal("expected missing extension to fail")
	}
}

func TestWithinRejectsParentDirectory(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(root, "..", "outside.py")
	if within(root, outside) {
		t.Fatal("expected path outside plugin directory to be rejected")
	}
}

func TestRunRejectsSymlinkEscape(t *testing.T) {
	directory := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.py")
	if err := os.WriteFile(outside, []byte("pass\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(directory, "escape.py")); err != nil {
		t.Skipf("symlinks are unavailable: %v", err)
	}
	err := Run(bytes.NewReader(nil), &bytes.Buffer{}, &bytes.Buffer{}, directory, "escape")
	if err == nil {
		t.Fatal("expected symlink escaping plugin directory to be rejected")
	}
}
