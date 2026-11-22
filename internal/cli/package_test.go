package cli

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestPackageAppliesBPackIgnoreRules(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"Bitfile":               "blueprint \"demo\" v1\nbring \"./\"\n",
		".bpackignore":          "node_modules/\n.env\n*.log\n!logs/keep.log\n",
		"src/main.go":           "package main\n",
		"node_modules/pkg/a.js": "ignored",
		".env":                  "SECRET=value",
		"logs/keep.log":         "keep",
		"logs/drop.log":         "ignore",
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	outputPath := filepath.Join(t.TempDir(), "demo.bpkg")
	if err := bitpackPackage(io.Discard, []string{filepath.Join(root, "Bitfile"), outputPath}); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	compressed, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer compressed.Close()
	archive := tar.NewReader(compressed)
	included := make(map[string]bool)
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		included[header.Name] = true
	}
	for _, name := range []string{"Bitfile", ".bpackignore", "src/main.go", "logs/keep.log"} {
		if !included[name] {
			t.Errorf("expected %q in package", name)
		}
	}
	for _, name := range []string{"node_modules/pkg/a.js", ".env", "logs/drop.log"} {
		if included[name] {
			t.Errorf("unexpected ignored file %q in package", name)
		}
	}
}

func TestPackageAndUnpackRoundTripProjectFiles(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"Bitfile":      "blueprint \"demo\" v1\nbring \"./\"\n",
		".env":         "API_KEY=private",
		".bpackignore": "dist/\n",
		"src/main.go":  "package main\n",
		"data/sample":  "application data\n",
		"dist/bundle":  "generated",
	}
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	packagePath := filepath.Join(t.TempDir(), "demo.bpkg")
	if err := bitpackPackage(io.Discard, []string{filepath.Join(root, "Bitfile"), packagePath}); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "unpacked")
	if err := bitpackUnpack(io.Discard, []string{packagePath, destination}); err != nil {
		t.Fatal(err)
	}
	for name, expected := range map[string]string{
		"src/main.go": "package main\n",
		"data/sample": "application data\n",
	} {
		content, err := os.ReadFile(filepath.Join(destination, filepath.FromSlash(name)))
		if err != nil {
			t.Errorf("read extracted file %s: %v", name, err)
			continue
		}
		if string(content) != expected {
			t.Errorf("unexpected extracted content for %s: %q", name, content)
		}
	}
	for _, name := range []string{".env", "dist/bundle"} {
		if _, err := os.Stat(filepath.Join(destination, filepath.FromSlash(name))); !os.IsNotExist(err) {
			t.Errorf("expected ignored file %s to be absent, stat error: %v", name, err)
		}
	}
}
