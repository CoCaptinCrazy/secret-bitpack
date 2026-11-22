package cli

import (
	"fmt"
	"path/filepath"
)

var dataDirectory = "data"

func setDataDirectory(path string) error {
	if path == "" {
		return fmt.Errorf("data directory cannot be empty")
	}
	dataDirectory = filepath.Clean(path)
	return nil
}

func dataPath(parts ...string) string {
	all := append([]string{dataDirectory}, parts...)
	return filepath.Join(all...)
}
