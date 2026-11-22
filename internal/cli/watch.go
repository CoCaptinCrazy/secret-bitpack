package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func bitpackWatch(output io.Writer, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: watch [directory]")
	}
	root := "."
	if len(args) == 1 {
		root = args[0]
	}
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("watch %s: %w", root, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("watch target %s is not a directory", root)
	}
	previous, err := snapshotDirectory(root)
	if err != nil {
		return err
	}
	fmt.Fprintf(output, "Watching %s; press Ctrl+C to stop\n", root)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		current, err := snapshotDirectory(root)
		if err != nil {
			return err
		}
		for path, modified := range current {
			old, exists := previous[path]
			switch {
			case !exists:
				fmt.Fprintf(output, "[CREATE] %s\n", path)
			case !old.Equal(modified):
				fmt.Fprintf(output, "[CHANGE] %s\n", path)
			}
		}
		for path := range previous {
			if _, exists := current[path]; !exists {
				fmt.Fprintf(output, "[DELETE] %s\n", path)
			}
		}
		previous = current
	}
	return nil
}

func snapshotDirectory(root string) (map[string]time.Time, error) {
	result := make(map[string]time.Time)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != root && (entry.Name() == ".git" || entry.Name() == "node_modules" || entry.Name() == ".bitpack") {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		result[path] = info.ModTime()
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan watch directory: %w", err)
	}
	return result, nil
}
