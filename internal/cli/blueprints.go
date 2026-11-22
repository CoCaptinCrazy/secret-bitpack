package cli

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/builder/blueprints"
	"github.com/CoCaptinCrazy/secret-bitpack/internal/builder/parser"
	devruntime "github.com/CoCaptinCrazy/secret-bitpack/internal/dev/runtime"
	bitpacksecrets "github.com/CoCaptinCrazy/secret-bitpack/internal/security/secrets"
)

var validResourceName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`)

func bitpackForge(output io.Writer, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: forge [Bitfile]")
	}
	path := "Bitfile"
	if len(args) == 1 {
		path = args[0]
	}
	source, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read Bitfile %s: %w", path, err)
	}
	parsed, err := parser.Parse(strings.NewReader(string(source)))
	if err != nil {
		return fmt.Errorf("parse Bitfile %s: %w", path, err)
	}
	if !validResourceName.MatchString(parsed.Name) {
		return fmt.Errorf("blueprint name %q contains unsupported characters", parsed.Name)
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve Bitfile path %s: %w", path, err)
	}
	workspace, err := resolveBlueprintWorkspace(filepath.Dir(absolutePath), parsed.Context)
	if err != nil {
		return err
	}

	id := sha256.Sum256(source)
	now := time.Now().Unix()
	record := map[string]any{
		"kind":       "Blueprint",
		"apiVersion": "v1",
		"metadata": map[string]any{
			"id": idHex(id[:]), "name": parsed.Name, "version": parsed.Version,
			"created": now, "author": "unknown", "description": "",
		},
		"runtime": map[string]any{
			"using": parsed.Using, "environment": parsed.Environment, "health": parsed.Health,
		},
		"execution": map[string]any{
			"assemble": parsed.Assemble, "verify": parsed.Verify, "ignite": parsed.Ignite,
		},
		"storage":    map[string]any{"volumes": []string{}, "caches": []string{}, "snapshots": true},
		"networking": map[string]any{"meshes": []string{}, "bridges": []string{}, "publishedPorts": []int{}},
		"security":   map[string]any{"secrets": parsed.Secrets, "signatures": []string{}, "policies": []string{}},
		"monitoring": map[string]any{"metrics": true, "alerts": true, "profiling": true},
		"assembly":   map[string]any{"compatible": true, "clusterReady": true, "autoscaleReady": true},
		"layers":     []string{},
		"history":    []map[string]any{{"action": "forge", "timestamp": now}},
		"source": map[string]any{
			"bitfile": string(source), "context": parsed.Context,
			"workspace": workspace, "sections": parsed.Sections, "directives": parsed.Directives,
		},
	}
	if err := os.MkdirAll(dataPath("blueprints"), 0o700); err != nil {
		return fmt.Errorf("create blueprint storage: %w", err)
	}
	destination := dataPath("blueprints", parsed.Name+".bp")
	if err := writeJSON(destination, record, 0o600); err != nil {
		return fmt.Errorf("write blueprint: %w", err)
	}
	fmt.Fprintf(output, "Blueprint created: %s\n", destination)
	return nil
}

func bitpackVerify(output io.Writer, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: verify [blueprint-or-Bitfile]")
	}
	target := "Bitfile"
	if len(args) == 1 {
		target = args[0]
	}
	parsed, workspace, document, err := loadBlueprintInput(target)
	if err != nil {
		return err
	}
	if document != nil {
		for _, field := range []string{"metadata", "runtime", "execution"} {
			if _, ok := document[field]; !ok {
				return fmt.Errorf("blueprint %q is invalid: missing %s", parsed.Name, field)
			}
		}
	}
	if len(parsed.Verify) == 0 {
		fmt.Fprintln(output, "No verify commands defined.")
	} else {
		environment, err := blueprintEnvironment(parsed)
		if err != nil {
			return err
		}
		if err := devruntime.RunCommands(
			context.Background(), parsed.Verify, environment, workspace, os.Stdin, output, os.Stderr,
		); err != nil {
			return fmt.Errorf("verify blueprint %q: %w", parsed.Name, err)
		}
	}
	fmt.Fprintf(output, "Blueprint is valid: %s\n", parsed.Name)
	return nil
}

func blueprintEnvironment(blueprint *blueprints.Blueprint) (map[string]string, error) {
	environment := make(map[string]string, len(blueprint.Environment)+len(blueprint.Secrets))
	for key, value := range blueprint.Environment {
		environment[key] = value
	}
	if len(blueprint.Secrets) == 0 {
		return environment, nil
	}
	values, err := bitpacksecrets.ResolveFromEnvironment(blueprint.Secrets)
	if err != nil {
		return nil, err
	}
	for name, value := range values {
		key := bitpacksecrets.EnvironmentKey(name)
		if _, exists := environment[key]; exists {
			return nil, fmt.Errorf("secret %q conflicts with environment variable %q", name, key)
		}
		environment[key] = value
	}
	return environment, nil
}

func loadBlueprintInput(target string) (*blueprints.Blueprint, string, map[string]any, error) {
	info, statErr := os.Stat(target)
	if statErr == nil && !info.IsDir() {
		if filepath.Ext(target) == ".bp" {
			document, err := readBlueprint(target)
			if err != nil {
				return nil, "", nil, err
			}
			parsed, workspace, err := parsedBlueprintFromDocument(document)
			return parsed, workspace, document, err
		}
		parsed, workspace, err := parseBlueprintFile(target)
		return parsed, workspace, nil, err
	}
	if statErr != nil && !os.IsNotExist(statErr) {
		return nil, "", nil, fmt.Errorf("inspect blueprint input %s: %w", target, statErr)
	}
	if !validResourceName.MatchString(target) {
		return nil, "", nil, fmt.Errorf("invalid blueprint name %q", target)
	}
	document, err := readBlueprint(dataPath("blueprints", target+".bp"))
	if err != nil {
		return nil, "", nil, err
	}
	parsed, workspace, err := parsedBlueprintFromDocument(document)
	return parsed, workspace, document, err
}

func parseBlueprintFile(path string) (*blueprints.Blueprint, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("open blueprint source %s: %w", path, err)
	}
	defer file.Close()
	parsed, err := parser.Parse(file)
	if err != nil {
		return nil, "", fmt.Errorf("invalid blueprint source %s: %w", path, err)
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("resolve blueprint source %s: %w", path, err)
	}
	workspace, err := resolveBlueprintWorkspace(filepath.Dir(absolutePath), parsed.Context)
	if err != nil {
		return nil, "", err
	}
	return parsed, workspace, nil
}

func parsedBlueprintFromDocument(document map[string]any) (*blueprints.Blueprint, string, error) {
	source, ok := document["source"].(map[string]any)
	if !ok {
		return nil, "", fmt.Errorf("blueprint has invalid source metadata")
	}
	bitfile, ok := source["bitfile"].(string)
	if !ok {
		return nil, "", fmt.Errorf("blueprint has no Bitfile source")
	}
	parsed, err := parser.Parse(strings.NewReader(bitfile))
	if err != nil {
		return nil, "", fmt.Errorf("parse stored Bitfile: %w", err)
	}
	if workspace, ok := source["workspace"].(string); ok && workspace != "" {
		return parsed, workspace, nil
	}
	workspace, err := resolveBlueprintWorkspace(".", parsed.Context)
	if err != nil {
		return nil, "", err
	}
	return parsed, workspace, nil
}

func resolveBlueprintWorkspace(base, contextPath string) (string, error) {
	if contextPath == "" {
		contextPath = "."
	}
	if !filepath.IsAbs(contextPath) {
		contextPath = filepath.Join(base, contextPath)
	}
	workspace, err := filepath.Abs(contextPath)
	if err != nil {
		return "", fmt.Errorf("resolve blueprint workspace: %w", err)
	}
	return workspace, nil
}

func bitpackBlueprints(output io.Writer) error {
	entries, err := os.ReadDir(dataPath("blueprints"))
	if os.IsNotExist(err) {
		fmt.Fprintln(output, "No blueprints found.")
		return nil
	}
	if err != nil {
		return fmt.Errorf("list blueprints: %w", err)
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".bp" {
			continue
		}
		document, err := readBlueprint(dataPath("blueprints", entry.Name()))
		if err != nil {
			return err
		}
		metadata, ok := document["metadata"].(map[string]any)
		if !ok {
			return fmt.Errorf("blueprint %s has invalid metadata", entry.Name())
		}
		fmt.Fprintf(output, "%s\t%s\t%s\n",
			valueString(metadata["name"]), valueString(metadata["version"]), valueString(metadata["id"]))
		count++
	}
	if count == 0 {
		fmt.Fprintln(output, "No blueprints found.")
	}
	return nil
}

func bitpackInspect(output io.Writer, args []string) error {
	if len(args) != 1 || !validResourceName.MatchString(args[0]) {
		return fmt.Errorf("usage: inspect <blueprint>")
	}
	document, err := readBlueprint(dataPath("blueprints", args[0]+".bp"))
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(document)
}

func bitpackExport(output io.Writer, args []string) error {
	if len(args) != 1 || !validResourceName.MatchString(args[0]) {
		return fmt.Errorf("usage: export <blueprint>")
	}
	source := dataPath("blueprints", args[0]+".bp")
	if _, err := readBlueprint(source); err != nil {
		return err
	}
	target := args[0] + ".export.bp"
	if err := copyFile(source, target, 0o600); err != nil {
		return fmt.Errorf("export blueprint: %w", err)
	}
	fmt.Fprintf(output, "Exported to %s\n", target)
	return nil
}

func bitpackImport(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: import <file.bp>")
	}
	source := args[0]
	if filepath.Ext(source) != ".bp" {
		return fmt.Errorf("import source must have a .bp extension")
	}
	document, err := readBlueprint(source)
	if err != nil {
		return err
	}
	metadata, ok := document["metadata"].(map[string]any)
	if !ok {
		return fmt.Errorf("imported blueprint has no metadata")
	}
	name := valueString(metadata["name"])
	if !validResourceName.MatchString(name) {
		return fmt.Errorf("imported blueprint has invalid name %q", name)
	}
	destination := dataPath("blueprints", name+".bp")
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return fmt.Errorf("create blueprint storage: %w", err)
	}
	if err := copyFile(source, destination, 0o600); err != nil {
		return fmt.Errorf("import blueprint: %w", err)
	}
	fmt.Fprintf(output, "Imported: %s\n", destination)
	return nil
}

func readBlueprint(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read blueprint %s: %w", path, err)
	}
	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("decode blueprint %s: %w", path, err)
	}
	return document, nil
}

func writeJSON(path string, value any, mode os.FileMode) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, data, mode)
}

func copyFile(source, destination string, mode os.FileMode) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return os.WriteFile(destination, data, mode)
}

func valueString(value any) string {
	text, _ := value.(string)
	if text == "" {
		return "unknown"
	}
	return text
}

func idHex(value []byte) string {
	return hex.EncodeToString(value)
}
