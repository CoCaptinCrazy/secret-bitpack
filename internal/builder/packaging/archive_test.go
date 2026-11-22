package packaging

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteGzipTarEntriesSortsAndPreservesModes(t *testing.T) {
	var archive bytes.Buffer
	err := WriteGzipTarEntries(&archive, []Entry{
		{Path: "src/tool.sh", Content: []byte("run"), Mode: 0o755},
		{Path: "Bitfile", Content: []byte("blueprint")},
	})
	if err != nil {
		t.Fatal(err)
	}
	compressed, err := gzip.NewReader(bytes.NewReader(archive.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	defer compressed.Close()
	reader := tar.NewReader(compressed)
	first, err := reader.Next()
	if err != nil {
		t.Fatal(err)
	}
	if first.Name != "Bitfile" {
		t.Fatalf("archive entries are not sorted: %q", first.Name)
	}
	second, err := reader.Next()
	if err != nil {
		t.Fatal(err)
	}
	if second.Name != "src/tool.sh" || second.Mode != 0o755 {
		t.Fatalf("mode or path was not preserved: %#v", second)
	}
	_, err = reader.Next()
	if err != io.EOF {
		t.Fatalf("expected end of archive, got %v", err)
	}
}

func TestWriteGzipTarEntriesRejectsTraversal(t *testing.T) {
	var archive bytes.Buffer
	if err := WriteGzipTarEntries(&archive, []Entry{{Path: "../outside", Content: []byte("bad")}}); err == nil {
		t.Fatal("expected archive traversal path to fail")
	}
}

func TestExtractGzipTarRestoresFilesSafely(t *testing.T) {
	var archive bytes.Buffer
	if err := WriteGzipTarEntries(&archive, []Entry{
		{Path: "src/main.go", Content: []byte("package main\n"), Mode: 0o644},
	}); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "extracted")
	if err := ExtractGzipTar(bytes.NewReader(archive.Bytes()), destination); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(destination, "src", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "package main\n" {
		t.Fatalf("unexpected extracted content %q", content)
	}
}

func TestExtractGzipTarRejectsCorruptGzipChecksum(t *testing.T) {
	var archive bytes.Buffer
	if err := WriteGzipTarEntries(&archive, []Entry{{Path: "file.txt", Content: []byte("valid")}}); err != nil {
		t.Fatal(err)
	}
	corrupt := append([]byte(nil), archive.Bytes()...)
	corrupt[len(corrupt)-8] ^= 0xff
	if err := ExtractGzipTar(bytes.NewReader(corrupt), filepath.Join(t.TempDir(), "extracted")); err == nil {
		t.Fatal("expected corrupt gzip checksum to be rejected")
	}
}

func TestExtractGzipTarRejectsTraversalAndCleansDestination(t *testing.T) {
	var archive bytes.Buffer
	compressed := gzip.NewWriter(&archive)
	writer := tar.NewWriter(compressed)
	if err := writer.WriteHeader(&tar.Header{Name: "../outside.txt", Mode: 0o644, Size: 4}); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("evil")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := compressed.Close(); err != nil {
		t.Fatal(err)
	}

	parent := t.TempDir()
	destination := filepath.Join(parent, "extracted")
	if err := ExtractGzipTar(bytes.NewReader(archive.Bytes()), destination); err == nil {
		t.Fatal("expected traversal entry to be rejected")
	}
	if _, err := os.Lstat(destination); !os.IsNotExist(err) {
		t.Fatalf("expected failed extraction to be cleaned up, got %v", err)
	}
	if _, err := os.Lstat(filepath.Join(parent, "outside.txt")); !os.IsNotExist(err) {
		t.Fatalf("archive wrote outside the extraction root: %v", err)
	}
}
