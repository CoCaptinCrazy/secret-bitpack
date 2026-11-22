package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Change struct {
	Path    string
	ModTime time.Time
}

type FileStamp struct {
	ModTime time.Time
	Size    int64
}

func Changed(previous, current Change) bool {
	return previous.Path != current.Path || !previous.ModTime.Equal(current.ModTime)
}

func Snapshot(root string, paths []string) (map[string]FileStamp, error) {
	result := make(map[string]FileStamp)
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve watch root: %w", err)
	}
	for _, path := range paths {
		if filepath.IsAbs(path) {
			return nil, fmt.Errorf("watch path %q must be relative to the workspace", path)
		}
		target := filepath.Join(absoluteRoot, filepath.FromSlash(path))
		relative, err := filepath.Rel(absoluteRoot, target)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("watch path %q escapes the workspace", path)
		}
		if err := snapshotPath(absoluteRoot, target, result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func ChangedPaths(previous, current map[string]FileStamp) []string {
	changed := make([]string, 0)
	for path, stamp := range current {
		old, exists := previous[path]
		if !exists || old != stamp {
			changed = append(changed, path)
		}
	}
	for path := range previous {
		if _, exists := current[path]; !exists {
			changed = append(changed, path)
		}
	}
	return changed
}

func snapshotPath(root, target string, result map[string]FileStamp) error {
	info, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect watch path %s: %w", target, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	if !info.IsDir() {
		result[target] = FileStamp{ModTime: info.ModTime(), Size: info.Size()}
		return nil
	}
	return filepath.WalkDir(target, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path != target && entry.IsDir() &&
			(entry.Name() == ".git" || entry.Name() == ".bitpack" ||
				entry.Name() == "node_modules" || entry.Name() == "data") {
			return filepath.SkipDir
		}
		if entry.Type()&os.ModeSymlink != 0 {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		fileInfo, err := entry.Info()
		if err != nil {
			return err
		}
		result[path] = FileStamp{ModTime: fileInfo.ModTime(), Size: fileInfo.Size()}
		return nil
	})
}
