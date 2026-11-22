package packaging

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Entry struct {
	Path    string
	Content []byte
	Mode    int64
}

func WriteGzipTar(destination io.Writer, files map[string][]byte) error {
	entries := make([]Entry, 0, len(files))
	for name, content := range files {
		entries = append(entries, Entry{Path: name, Content: content, Mode: 0o644})
	}
	return WriteGzipTarEntries(destination, entries)
}

func WriteGzipTarEntries(destination io.Writer, entries []Entry) error {
	gzipWriter := gzip.NewWriter(destination)
	gzipWriter.ModTime = time.Unix(0, 0).UTC()
	if err := WriteTarEntries(gzipWriter, entries); err != nil {
		_ = gzipWriter.Close()
		return err
	}
	return gzipWriter.Close()
}

func WriteTarEntries(destination io.Writer, entries []Entry) error {
	ordered := append([]Entry(nil), entries...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Path < ordered[j].Path })
	tarWriter := tar.NewWriter(destination)
	for _, entry := range ordered {
		name, err := safeArchivePath(entry.Path)
		if err != nil {
			_ = tarWriter.Close()
			return err
		}
		mode := entry.Mode & 0o777
		if mode == 0 {
			mode = 0o644
		}
		header := &tar.Header{
			Name: name, Mode: mode, Size: int64(len(entry.Content)),
			ModTime: time.Unix(0, 0).UTC(), Typeflag: tar.TypeReg,
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			_ = tarWriter.Close()
			return err
		}
		if _, err := tarWriter.Write(entry.Content); err != nil {
			_ = tarWriter.Close()
			return err
		}
	}
	return tarWriter.Close()
}

func safeArchivePath(name string) (string, error) {
	if name == "" || filepath.IsAbs(name) {
		return "", fmt.Errorf("archive path %q must be relative", name)
	}
	clean := filepath.Clean(filepath.FromSlash(name))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("archive path %q escapes the project root", name)
	}
	return filepath.ToSlash(clean), nil
}

func ExtractGzipTar(source io.Reader, destination string) (returnErr error) {
	root, err := filepath.Abs(destination)
	if err != nil {
		return fmt.Errorf("resolve extraction destination: %w", err)
	}
	if _, err := os.Lstat(root); err == nil {
		return fmt.Errorf("extraction destination already exists: %s", root)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect extraction destination: %w", err)
	}
	if err := os.Mkdir(root, 0o700); err != nil {
		return fmt.Errorf("create extraction destination: %w", err)
	}
	completed := false
	defer func() {
		if !completed {
			if cleanupErr := os.RemoveAll(root); cleanupErr != nil {
				if returnErr == nil {
					returnErr = fmt.Errorf("remove incomplete package extraction: %w", cleanupErr)
				} else {
					returnErr = fmt.Errorf("%w; also failed to remove incomplete extraction: %v", returnErr, cleanupErr)
				}
			}
		}
	}()

	compressed, err := gzip.NewReader(source)
	if err != nil {
		return fmt.Errorf("open compressed package: %w", err)
	}
	closed := false
	defer func() {
		if !closed {
			if err := compressed.Close(); err != nil && returnErr == nil {
				returnErr = fmt.Errorf("close compressed package: %w", err)
			}
		}
	}()
	reader := tar.NewReader(compressed)
	totalBytes := int64(0)
	fileCount := 0
	entryCount := 0
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read package entry: %w", err)
		}
		entryCount++
		if entryCount > 100000 {
			return fmt.Errorf("package exceeds extraction entry limit")
		}
		name, err := safeArchivePath(header.Name)
		if err != nil {
			return err
		}
		target := filepath.Join(root, filepath.FromSlash(name))
		if relative, err := filepath.Rel(root, target); err != nil ||
			relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return fmt.Errorf("package entry %q escapes extraction destination", header.Name)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o700); err != nil {
				return fmt.Errorf("create package directory %s: %w", header.Name, err)
			}
		case tar.TypeReg, tar.TypeRegA:
			fileCount++
			if fileCount > 100000 || header.Size < 0 || header.Size > (10<<30)-totalBytes {
				return fmt.Errorf("package exceeds extraction limits")
			}
			totalBytes += header.Size
			if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
				return fmt.Errorf("create package directory for %s: %w", header.Name, err)
			}
			mode := os.FileMode(header.Mode).Perm() & 0o755
			if mode == 0 {
				mode = 0o644
			}
			file, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
			if err != nil {
				return fmt.Errorf("create package file %s: %w", header.Name, err)
			}
			_, copyErr := io.CopyN(file, reader, header.Size)
			closeErr := file.Close()
			if copyErr != nil {
				return fmt.Errorf("extract package file %s: %w", header.Name, copyErr)
			}
			if closeErr != nil {
				return fmt.Errorf("close package file %s: %w", header.Name, closeErr)
			}
		default:
			return fmt.Errorf("package entry %q has unsupported type %d", header.Name, header.Typeflag)
		}
	}
	trailingBytes, err := io.CopyN(io.Discard, compressed, (1<<20)+1)
	if err != nil && err != io.EOF {
		return fmt.Errorf("validate compressed package: %w", err)
	}
	if trailingBytes > 1<<20 {
		return fmt.Errorf("package has excessive trailing data")
	}
	if err := compressed.Close(); err != nil {
		return fmt.Errorf("close compressed package: %w", err)
	}
	closed = true
	completed = true
	return nil
}
