package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/builder/packaging"
)

func bitpackPackage(output io.Writer, args []string) error {
	if len(args) > 2 {
		return fmt.Errorf("usage: pack [Bitfile] [output.bpkg]")
	}
	source := "Bitfile"
	if len(args) > 0 {
		source = args[0]
	}
	parsed, workspace, err := parseBlueprintFile(source)
	if err != nil {
		return err
	}
	destination := parsed.Name + ".bpkg"
	if len(args) == 2 {
		destination = args[1]
	}
	absoluteDestination, err := filepath.Abs(destination)
	if err != nil {
		return fmt.Errorf("resolve package destination: %w", err)
	}
	entries, err := packaging.CollectDirectory(workspace, absoluteDestination)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("project context %s contains no packageable files", workspace)
	}
	if err := os.MkdirAll(filepath.Dir(absoluteDestination), 0o700); err != nil {
		return fmt.Errorf("create package output directory: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(absoluteDestination), ".bitpack-package-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary package: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := packaging.WriteGzipTarEntries(temporary, entries); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("write project package: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close project package: %w", err)
	}
	if err := os.Rename(temporaryPath, absoluteDestination); err != nil {
		return fmt.Errorf("finalize project package %s: %w", absoluteDestination, err)
	}
	fmt.Fprintf(output, "Package created: %s (%d files)\n", absoluteDestination, len(entries))
	return nil
}

func bitpackUnpack(output io.Writer, args []string) error {
	if len(args) < 1 || len(args) > 2 || filepath.Ext(args[0]) != ".bpkg" {
		return fmt.Errorf("usage: unpack <file.bpkg> [destination]")
	}
	destination := strings.TrimSuffix(args[0], ".bpkg")
	if len(args) == 2 {
		destination = args[1]
	}
	file, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open Bitpack package %s: %w", args[0], err)
	}
	if err := packaging.ExtractGzipTar(file, destination); err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return fmt.Errorf("unpack Bitpack package: %w; close package: %v", err, closeErr)
		}
		return fmt.Errorf("unpack Bitpack package: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close Bitpack package: %w", err)
	}
	fmt.Fprintf(output, "Package extracted to %s\n", destination)
	return nil
}
