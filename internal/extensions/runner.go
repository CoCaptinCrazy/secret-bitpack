package extensions

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var validName = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func Run(input io.Reader, output, errorOutput io.Writer, directory, name string, args ...string) error {
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid extension name %q; use letters, numbers, dashes, or underscores", name)
	}
	root, err := filepath.Abs(directory)
	if err != nil {
		return fmt.Errorf("resolve extensions directory: %w", err)
	}
	script, err := filepath.Abs(filepath.Join(root, name+".py"))
	if err != nil {
		return fmt.Errorf("resolve extension path: %w", err)
	}
	if !within(root, script) {
		return fmt.Errorf("extension %q is outside the extensions directory", name)
	}
	info, err := os.Stat(script)
	if err != nil {
		return fmt.Errorf("find Python extension %q in %s: %w", name, root, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("Python extension %q is not a regular file", name)
	}
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return fmt.Errorf("resolve extensions directory: %w", err)
	}
	resolvedScript, err := filepath.EvalSymlinks(script)
	if err != nil {
		return fmt.Errorf("resolve extension path: %w", err)
	}
	if !within(resolvedRoot, resolvedScript) {
		return fmt.Errorf("Python extension %q resolves outside the extensions directory", name)
	}

	command, err := pythonCommand(script, args)
	if err != nil {
		return err
	}
	command.Stdin = input
	command.Stdout = output
	command.Stderr = errorOutput
	command.Env = append(os.Environ(),
		"BITPACK_HOST=go",
		"BITPACK_EXTENSION="+name,
	)
	if err := command.Run(); err != nil {
		return fmt.Errorf("Python extension %q failed: %w", name, err)
	}
	return nil
}

func pythonCommand(script string, args []string) (*exec.Cmd, error) {
	if executable := os.Getenv("BITPACK_PYTHON"); executable != "" {
		path, err := exec.LookPath(executable)
		if err != nil {
			return nil, fmt.Errorf("BITPACK_PYTHON=%q is unavailable: %w", executable, err)
		}
		return exec.CommandContext(context.Background(), path, append([]string{script}, args...)...), nil
	}
	if executable, err := exec.LookPath("python"); err == nil {
		return exec.CommandContext(context.Background(), executable, append([]string{script}, args...)...), nil
	}
	if executable, err := exec.LookPath("py"); err == nil {
		return exec.CommandContext(context.Background(), executable, append([]string{"-3", script}, args...)...), nil
	}
	return nil, fmt.Errorf("Python was not found; install Python or set BITPACK_PYTHON to its executable")
}

func within(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) &&
		!filepath.IsAbs(relative)
}
