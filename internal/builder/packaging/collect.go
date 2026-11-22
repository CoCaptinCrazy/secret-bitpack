package packaging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"secret-bitpack/internal/ignore"
)

func CollectDirectory(workspace, excludedPath string) ([]Entry, error) {
	root, err := filepath.Abs(workspace)
	if err != nil {
		return nil, fmt.Errorf("resolve project workspace: %w", err)
	}
	excluded := ""
	if excludedPath != "" {
		excluded, err = filepath.Abs(excludedPath)
		if err != nil {
			return nil, fmt.Errorf("resolve excluded package path: %w", err)
		}
	}
	matcher, err := ignore.Load(filepath.Join(root, ".bpackignore"))
	if err != nil {
		return nil, err
	}
	entries := make([]Entry, 0)
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("resolve relative package path %s: %w", path, err)
		}
		relativePath := filepath.ToSlash(relative)
		if entry.IsDir() {
			if matcher.Ignored(relativePath, true) {
				return filepath.SkipDir
			}
			return nil
		}
		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("resolve package input %s: %w", path, err)
		}
		samePath := absolutePath == excluded
		if !samePath && excluded != "" && os.PathSeparator == '\\' {
			samePath = strings.EqualFold(absolutePath, excluded)
		}
		if samePath || matcher.Ignored(relativePath, false) {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect package input %s: %w", path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to package symlink %s", path)
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read package input %s: %w", path, err)
		}
		entries = append(entries, Entry{
			Path: relativePath, Content: content, Mode: int64(info.Mode().Perm()),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect project package: %w", err)
	}
	return entries, nil
}
