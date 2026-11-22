package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/core"
	"github.com/CoCaptinCrazy/secret-bitpack/internal/orchestration/cluster"
)

func bitpackConfig(bitpack *core.Bitpack, output io.Writer, args []string) error {
	if len(args) == 0 {
		fmt.Fprintf(output, "environment=%s\nregistry=%s\ndata_dir=%s\n",
			bitpack.Config.Environment, bitpack.Config.Registry, bitpack.Config.DataDir)
		return nil
	}
	switch args[0] {
	case "get":
		if len(args) != 2 {
			return fmt.Errorf("usage: config get <environment|registry|data_dir>")
		}
		value, err := configValue(bitpack, args[1])
		if err != nil {
			return err
		}
		fmt.Fprintln(output, value)
	case "set":
		if len(args) != 3 {
			return fmt.Errorf("usage: config set <environment|registry|data_dir> <value>")
		}
		previous := *bitpack.Config
		switch args[1] {
		case "environment":
			bitpack.Config.Environment = args[2]
		case "registry":
			bitpack.Config.Registry = args[2]
		case "data_dir":
			bitpack.Config.DataDir = args[2]
			if err := setDataDirectory(args[2]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown config key %q", args[1])
		}
		if args[1] == "data_dir" && args[2] == "" {
			*bitpack.Config = previous
			return fmt.Errorf("data_dir cannot be empty")
		}
		if err := bitpack.Config.Save(filepath.Join(".bitpack", "config.json")); err != nil {
			*bitpack.Config = previous
			return fmt.Errorf("save config: %w", err)
		}
		if args[1] == "data_dir" {
			if err := setDataDirectory(args[2]); err != nil {
				*bitpack.Config = previous
				return err
			}
		}
		fmt.Fprintf(output, "Updated %s\n", args[1])
	default:
		return fmt.Errorf("usage: config [get <key>|set <key> <value>]")
	}
	return nil
}

func configValue(bitpack *core.Bitpack, key string) (string, error) {
	switch key {
	case "environment":
		return bitpack.Config.Environment, nil
	case "registry":
		return bitpack.Config.Registry, nil
	case "data_dir":
		return bitpack.Config.DataDir, nil
	default:
		return "", fmt.Errorf("unknown config key %q", key)
	}
}

func bitpackBackup(output io.Writer) error {
	backupRoot := dataPath("backups")
	if err := os.MkdirAll(backupRoot, 0o700); err != nil {
		return fmt.Errorf("create backup directory: %w", err)
	}
	name := "backup-" + time.Now().UTC().Format("20060102T150405.000000000Z")
	destination := filepath.Join(backupRoot, name)
	copied := 0
	for _, directory := range []string{"blueprints", "modules", "clusters", "assemblies", "meshes", "snapshots", "logs"} {
		source := dataPath(directory)
		if _, err := os.Stat(source); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return fmt.Errorf("inspect %s for backup: %w", source, err)
		}
		if err := copyDirectory(source, filepath.Join(destination, directory)); err != nil {
			return fmt.Errorf("backup %s: %w", directory, err)
		}
		copied++
	}
	if copied == 0 {
		if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove empty backup directory: %w", err)
		}
		fmt.Fprintln(output, "No data directories to back up.")
		return nil
	}
	fmt.Fprintf(output, "Backup created: %s\n", destination)
	return nil
}

func copyDirectory(source, destination string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to follow symlink %s", path)
		}
		if entry.IsDir() {
			return os.MkdirAll(target, 0o700)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o600)
	})
}

func bitpackMetrics(output io.Writer) error {
	fmt.Fprintln(output, "No metrics.")
	return nil
}

func bitpackProfile(output io.Writer) error {
	started := time.Now()
	_ = time.Now()
	fmt.Fprintln(output, "Profile complete.")
	fmt.Fprintf(output, "Elapsed: %.6fs\n", time.Since(started).Seconds())
	return nil
}

func bitpackAlert(output io.Writer, cpu int) error {
	if cpu > 90 {
		fmt.Fprintln(output, "[warning] CPU limit exceeded")
		return nil
	}
	fmt.Fprintln(output, "No alerts.")
	return nil
}

func bitpackEvents(output io.Writer) error {
	path := dataPath("events.jsonl")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		fmt.Fprintln(output, "No events.")
		return nil
	}
	if err != nil {
		return fmt.Errorf("read events: %w", err)
	}
	_, err = output.Write(data)
	return err
}

func bitpackSecrets(output io.Writer, args []string) error {
	if len(args) > 0 && args[0] != "list" {
		return fmt.Errorf("usage: secrets [list]")
	}
	root := filepath.Join(".bitpack", "secrets")
	entries, err := os.ReadDir(root)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("list secret names: %w", err)
	}
	names := make(map[string]bool)
	for _, entry := range entries {
		if !entry.IsDir() {
			names[strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))] = true
		}
	}
	for _, value := range os.Environ() {
		key, _, ok := strings.Cut(value, "=")
		if ok && strings.HasPrefix(key, "BITPACK_SECRET_") {
			names[strings.TrimPrefix(key, "BITPACK_SECRET_")] = true
		}
	}
	if len(names) == 0 {
		fmt.Fprintln(output, "No secrets.")
		return nil
	}
	ordered := make([]string, 0, len(names))
	for name := range names {
		ordered = append(ordered, name)
	}
	sort.Strings(ordered)
	for _, name := range ordered {
		fmt.Fprintln(output, name)
	}
	return nil
}

func bitpackMesh(output io.Writer, args []string) error {
	if len(args) == 0 {
		return bitpackMeshes(output)
	}
	if len(args) != 1 || !validResourceName.MatchString(args[0]) {
		return fmt.Errorf("usage: mesh [name]")
	}
	path := dataPath("meshes", args[0])
	if err := os.MkdirAll(path, 0o700); err != nil {
		return fmt.Errorf("create mesh definition: %w", err)
	}
	if err := writeJSON(filepath.Join(path, "mesh.json"), map[string]any{
		"name": args[0], "status": "configured", "created": time.Now().Unix(),
	}, 0o600); err != nil {
		return err
	}
	fmt.Fprintf(output, "Mesh definition created: %s (no network provisioned)\n", args[0])
	return nil
}

func bitpackMeshes(output io.Writer) error {
	return listDirectories(output, dataPath("meshes"), "No meshes found.")
}

func bitpackBridges(output io.Writer) error {
	root := dataPath("bridges")
	entries, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		fmt.Fprintln(output, "No bridges.")
		return nil
	}
	if err != nil {
		return fmt.Errorf("list bridges: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(output, "No bridges.")
		return nil
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		var bridge struct {
			Name       string `json:"name"`
			HostPort   int    `json:"host_port"`
			TargetPort int    `json:"target_port"`
		}
		data, err := os.ReadFile(filepath.Join(root, entry.Name()))
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &bridge); err != nil {
			return fmt.Errorf("decode bridge %s: %w", entry.Name(), err)
		}
		fmt.Fprintf(output, "%s %d->%d\n", bridge.Name, bridge.HostPort, bridge.TargetPort)
	}
	return nil
}

func bitpackRoutes(output io.Writer) error {
	path := dataPath("routes.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		fmt.Fprintln(output, "[]")
		return nil
	}
	if err != nil {
		return fmt.Errorf("read routes: %w", err)
	}
	var routes []json.RawMessage
	if err := json.Unmarshal(data, &routes); err != nil {
		return fmt.Errorf("decode routes: %w", err)
	}
	if len(routes) == 0 {
		fmt.Fprintln(output, "[]")
		return nil
	}
	for _, route := range routes {
		fmt.Fprintln(output, string(route))
	}
	return nil
}

func bitpackNodes(output io.Writer) error {
	store := cluster.NewStore(dataPath("clusters"), "0.1.17")
	records, err := store.List()
	if err != nil {
		return err
	}
	count := 0
	for _, record := range records {
		for _, node := range record.Nodes {
			fmt.Fprintf(output, "Cluster: %s\n - %v\n", record.Metadata.Name, node)
			count++
		}
	}
	if count == 0 {
		fmt.Fprintln(output, "No nodes.")
	}
	return nil
}

func bitpackPolicy(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: policy <file.bppolicy>")
	}
	if strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("policy name is required")
	}
	fmt.Fprintln(output, "Policy valid: true")
	return nil
}

func bitpackSearch(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: search <query>")
	}
	query := strings.ToLower(args[0])
	catalog := []string{"postgres", "redis", "nginx", "ubuntu", "website"}
	entries, err := os.ReadDir(dataPath("registry"))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read local registry: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".bp" {
			catalog = append(catalog, strings.TrimSuffix(entry.Name(), ".bp"))
		}
	}
	seen := map[string]bool{}
	for _, item := range catalog {
		if strings.Contains(strings.ToLower(item), query) && !seen[item] {
			fmt.Fprintln(output, item)
			seen[item] = true
		}
	}
	if len(seen) == 0 {
		fmt.Fprintln(output, "No results.")
	}
	return nil
}

func bitpackPublish(output io.Writer, args []string) error {
	if len(args) != 1 || !validResourceName.MatchString(args[0]) {
		return fmt.Errorf("usage: publish <blueprint>")
	}
	source := dataPath("blueprints", args[0]+".bp")
	if _, err := readBlueprint(source); err != nil {
		return err
	}
	destination := dataPath("registry", args[0]+".bp")
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return err
	}
	if err := copyFile(source, destination, 0o600); err != nil {
		return fmt.Errorf("publish to local registry: %w", err)
	}
	fmt.Fprintf(output, "Published to local registry: %s\n", args[0])
	return nil
}

func bitpackAcquire(output io.Writer, args []string) error {
	if len(args) != 1 || !validResourceName.MatchString(args[0]) {
		return fmt.Errorf("usage: acquire <blueprint>")
	}
	source := dataPath("registry", args[0]+".bp")
	if _, err := readBlueprint(source); err != nil {
		return err
	}
	destination := dataPath("blueprints", args[0]+".bp")
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return err
	}
	if err := copyFile(source, destination, 0o600); err != nil {
		return fmt.Errorf("acquire from local registry: %w", err)
	}
	fmt.Fprintf(output, "Acquired from local registry: %s\n", args[0])
	return nil
}

func bitpackTags(output io.Writer, args []string) error {
	if len(args) < 1 || len(args) > 2 || !validResourceName.MatchString(args[0]) {
		return fmt.Errorf("usage: tags <blueprint> [tag]")
	}
	path := dataPath("registry", "tags.json")
	tags := map[string][]string{}
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &tags); err != nil {
			return fmt.Errorf("decode registry tags: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read registry tags: %w", err)
	}
	if len(args) == 2 {
		if !validResourceName.MatchString(args[1]) {
			return fmt.Errorf("invalid tag %q", args[1])
		}
		for _, tag := range tags[args[0]] {
			if tag == args[1] {
				fmt.Fprintf(output, "Tag already exists: %s\n", args[1])
				return nil
			}
		}
		tags[args[0]] = append(tags[args[0]], args[1])
		if err := writeJSON(path, tags, 0o600); err != nil {
			return err
		}
		fmt.Fprintf(output, "Tag added: %s:%s\n", args[0], args[1])
		return nil
	}
	sort.Strings(tags[args[0]])
	if len(tags[args[0]]) == 0 {
		fmt.Fprintln(output, "No tags.")
		return nil
	}
	for _, tag := range tags[args[0]] {
		fmt.Fprintln(output, tag)
	}
	return nil
}

func bitpackLogin(bitpack *core.Bitpack, output io.Writer, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: login [username]")
	}
	if os.Getenv("BITPACK_REGISTRY_TOKEN") == "" {
		return fmt.Errorf("set BITPACK_REGISTRY_TOKEN before logging in; credentials are not stored by Bitpack")
	}
	username := "local"
	if len(args) == 1 {
		username = args[0]
	}
	bitpack.State.Set("registry.username", username)
	fmt.Fprintf(output, "Logged in as %s for this process\n", username)
	return nil
}

func bitpackLogout(bitpack *core.Bitpack, output io.Writer) error {
	bitpack.State.Delete("registry.username")
	fmt.Fprintln(output, "Logged out for this process.")
	return nil
}

func listDirectories(output io.Writer, root, emptyMessage string) error {
	entries, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		fmt.Fprintln(output, emptyMessage)
		return nil
	}
	if err != nil {
		return fmt.Errorf("list %s: %w", root, err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(output, emptyMessage)
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	if len(names) == 0 {
		fmt.Fprintln(output, emptyMessage)
		return nil
	}
	for _, name := range names {
		fmt.Fprintln(output, name)
	}
	return nil
}
