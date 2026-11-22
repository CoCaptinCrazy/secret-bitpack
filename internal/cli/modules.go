package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"secret-bitpack/internal/orchestration/kubernetes"
	"secret-bitpack/internal/runtime/modules"
	"secret-bitpack/internal/runtime/supervisor"
)

type snapshotRecord struct {
	ID      string         `json:"snapshot"`
	Module  string         `json:"module"`
	Created int64          `json:"created"`
	State   modules.Module `json:"state"`
}

func moduleStore() *modules.Store {
	return modules.NewStore(dataPath("modules"))
}

func processManager() *supervisor.Manager {
	return supervisor.New(dataPath("processes"), dataPath("logs", "modules"))
}

func findLocalModule(store *modules.Store, name string) (*modules.Module, error) {
	list, err := store.List()
	if err != nil {
		return nil, err
	}
	for index := range list {
		if list[index].Name == name {
			return &list[index], nil
		}
	}
	return nil, nil
}

func bitpackIgnite(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: ignite <blueprint>")
	}
	if definition, ok, err := readKubernetesAssembly(args[0]); err != nil {
		return err
	} else if ok {
		return restartKubernetesAssembly(output, definition)
	}
	parsed, workspace, document, err := loadBlueprintInput(args[0])
	if err != nil {
		return err
	}
	if len(parsed.Ignite) > 0 {
		if _, err := blueprintEnvironment(parsed); err != nil {
			return err
		}
		executable, commandArgs, err := supervisor.ShellCommand(strings.Join(parsed.Ignite, " && "))
		if err != nil {
			return err
		}
		store := moduleStore()
		module, err := findLocalModule(store, parsed.Name)
		if err != nil {
			return err
		}
		created := false
		if module == nil {
			module, err = store.Create(parsed.Name, parsed.Name, "local-process")
			if err != nil {
				return err
			}
			created = true
		} else if module.Runtime != "local-process" {
			return fmt.Errorf("module %q is already managed by %s", parsed.Name, module.Runtime)
		} else if running, err := processManager().IsRunning(module.ID); err == nil && running {
			return fmt.Errorf("module %q is already running; use restart to relaunch it", parsed.Name)
		} else if err != nil && !os.IsNotExist(err) {
			return err
		}
		spec := supervisor.Spec{
			Name: parsed.Name, Executable: executable, Arguments: commandArgs,
			WorkingDir: workspace, Environment: parsed.Environment, Secrets: parsed.Secrets,
		}
		record, err := processManager().Start(module.ID, spec)
		if err != nil {
			if created {
				if deleteErr := store.Delete(module); deleteErr != nil {
					return fmt.Errorf("start ignite process: %w; remove module state: %v", err, deleteErr)
				}
			}
			return fmt.Errorf("start ignite process for %q: %w", parsed.Name, err)
		}
		module.Runtime = "local-process"
		module.ProcessID = record.ID
		module.PID = record.PID
		module.Status = "running"
		module.Started = record.Started.Unix()
		if err := store.Update(module); err != nil {
			if stopErr := processManager().Stop(record.ID); stopErr != nil {
				return fmt.Errorf("save module state: %w; stop process: %v", err, stopErr)
			}
			return fmt.Errorf("save module state: %w", err)
		}
		fmt.Fprintf(output, "Background process started: %s (pid %d)\nLogs: %s\n",
			module.Name, record.PID, record.LogPath)
		return nil
	}
	if document == nil {
		fmt.Fprintf(output, "No ignite commands defined for %s.\n", parsed.Name)
		return nil
	}
	_, ok := document["metadata"].(map[string]any)
	if !ok {
		return fmt.Errorf("blueprint %q has invalid metadata", parsed.Name)
	}
	module, err := moduleStore().Create(
		parsed.Name,
		parsed.Name,
		parsed.Using,
	)
	if err != nil {
		return err
	}
	fmt.Fprintf(output, "Module state created: %s (%s); no runtime process launched\n",
		module.ID, module.Name)
	return nil
}

func bitpackModules(output io.Writer) error {
	list, err := moduleStore().List()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Fprintln(output, "No modules found.")
		return nil
	}
	fmt.Fprintf(output, "%-16s %-24s %-12s %-12s %s\n", "ID", "NAME", "STATUS", "RUNTIME", "PID")
	for _, module := range list {
		pid := ""
		if module.PID > 0 {
			pid = fmt.Sprintf("%d", module.PID)
		}
		fmt.Fprintf(output, "%-16s %-24s %-12s %-12s %s\n",
			module.ID, module.Name, module.Status, module.Runtime, pid)
	}
	return nil
}

func bitpackStop(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: stop <module-id-or-name>")
	}
	if definition, ok, err := readKubernetesAssembly(args[0]); err != nil {
		return err
	} else if ok {
		return stopKubernetesAssembly(output, definition)
	}
	store := moduleStore()
	module, err := store.Find(args[0])
	if err != nil {
		return err
	}
	if strings.HasPrefix(module.Runtime, "kubernetes:") {
		namespace := strings.TrimPrefix(module.Runtime, "kubernetes:")
		client := kubernetes.Client{Executable: os.Getenv("BITPACK_KUBECTL"), Namespace: namespace}
		if err := client.Stop(context.Background(), module.Name, output); err != nil {
			return err
		}
	} else if module.ProcessID != "" {
		if err := processManager().Stop(module.ProcessID); err != nil {
			return err
		}
	}
	module.Status = "stopped"
	if err := store.Update(module); err != nil {
		return err
	}
	fmt.Fprintf(output, "Module state stopped: %s\n", module.ID)
	return nil
}

func bitpackRestart(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: restart <module-id-or-name>")
	}
	if definition, ok, err := readKubernetesAssembly(args[0]); err != nil {
		return err
	} else if ok {
		return restartKubernetesAssembly(output, definition)
	}
	store := moduleStore()
	module, err := store.Find(args[0])
	if err != nil {
		return err
	}
	if strings.HasPrefix(module.Runtime, "kubernetes:") {
		namespace := strings.TrimPrefix(module.Runtime, "kubernetes:")
		client := kubernetes.Client{Executable: os.Getenv("BITPACK_KUBECTL"), Namespace: namespace}
		replicas := module.DesiredReplicas
		if replicas < 1 {
			replicas = 1
		}
		if err := client.Scale(context.Background(), module.Name, replicas, output); err != nil {
			return err
		}
		if err := client.Restart(context.Background(), module.Name, output); err != nil {
			return err
		}
		if module.AutoscaleMax > 0 {
			if err := client.ApplyAutoscaler(context.Background(), module.Blueprint, module.Name,
				module.AutoscaleMin, module.AutoscaleMax, module.AutoscaleCPU, output); err != nil {
				return err
			}
		}
		module.DesiredReplicas = replicas
	} else if module.ProcessID != "" {
		record, err := processManager().Restart(module.ProcessID)
		if err != nil {
			return err
		}
		module.Started = record.Started.Unix()
		module.PID = record.PID
	}
	module.Status = "running"
	module.Started = time.Now().Unix()
	if err := store.Update(module); err != nil {
		return err
	}
	if strings.HasPrefix(module.Runtime, "kubernetes:") {
		fmt.Fprintf(output, "Kubernetes module restarted: %s\n", module.ID)
	} else {
		fmt.Fprintf(output, "Module state restarted: %s (no runtime process launched)\n", module.ID)
	}
	return nil
}

func bitpackRemove(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: rm <module-id-or-name>")
	}
	store := moduleStore()
	module, err := store.Find(args[0])
	if err != nil {
		return err
	}
	if strings.HasPrefix(module.Runtime, "kubernetes:") {
		namespace := strings.TrimPrefix(module.Runtime, "kubernetes:")
		client := kubernetes.Client{Executable: os.Getenv("BITPACK_KUBECTL"), Namespace: namespace}
		if err := client.DeleteModule(context.Background(), module.Name, output); err != nil {
			return err
		}
	} else if module.ProcessID != "" {
		if err := processManager().Stop(module.ProcessID); err != nil {
			return err
		}
	}
	if err := store.Delete(module); err != nil {
		return err
	}
	fmt.Fprintf(output, "Module record removed: %s\n", module.ID)
	return nil
}

func bitpackInspectModule(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: inspect-module <module-id-or-name>")
	}
	module, err := moduleStore().Find(args[0])
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(module)
}

func bitpackRuntime(output io.Writer) error {
	fmt.Fprintln(output, "Runtime: Bitpack process supervisor and Kubernetes Deployments")
	return bitpackModules(output)
}

func bitpackHealth(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: health <module-id-or-name>")
	}
	if definition, ok, err := readKubernetesAssembly(args[0]); err != nil {
		return err
	} else if ok {
		client := kubernetes.Client{
			Executable: os.Getenv("BITPACK_KUBECTL"),
			Namespace:  definition.Namespace,
		}
		for _, resource := range definition.Resources {
			if resource.Kind == "module" {
				fmt.Fprintf(output, "%s: ", resource.Name)
				if err := client.Health(context.Background(), resource.Name, output); err != nil {
					return err
				}
			}
		}
		return nil
	}
	module, err := moduleStore().Find(args[0])
	if err != nil {
		return err
	}
	if strings.HasPrefix(module.Runtime, "kubernetes:") {
		namespace := strings.TrimPrefix(module.Runtime, "kubernetes:")
		client := kubernetes.Client{Executable: os.Getenv("BITPACK_KUBECTL"), Namespace: namespace}
		return client.Health(context.Background(), module.Name, output)
	}
	if module.ProcessID != "" {
		running, err := processManager().IsRunning(module.ProcessID)
		if err != nil {
			return err
		}
		if running {
			module.Status = "running"
		} else {
			module.Status = "stopped"
		}
		if record, err := processManager().Get(module.ProcessID); err != nil {
			return err
		} else {
			module.PID = record.PID
		}
		if err := moduleStore().Update(module); err != nil {
			return err
		}
		fmt.Fprintf(output, "Process health: %s (pid %d) for %s\n", module.Status, module.PID, module.ID)
		return nil
	}
	if module.Status != "running" {
		fmt.Fprintf(output, "Health status: stopped (metadata only) for %s\n", module.ID)
		return nil
	}
	fmt.Fprintf(output, "Health status: registered as running (process health is not monitored) for %s\n", module.ID)
	return nil
}

func bitpackLogs(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: logs <module-id-or-name>")
	}
	if definition, ok, err := readKubernetesAssembly(args[0]); err != nil {
		return err
	} else if ok {
		client := kubernetes.Client{
			Executable: os.Getenv("BITPACK_KUBECTL"),
			Namespace:  definition.Namespace,
		}
		for _, resource := range definition.Resources {
			if resource.Kind == "module" {
				fmt.Fprintf(output, "== %s ==\n", resource.Name)
				if err := client.Logs(context.Background(), resource.Name, output); err != nil {
					return err
				}
			}
		}
		return nil
	}
	module, err := moduleStore().Find(args[0])
	if err != nil {
		return err
	}
	if strings.HasPrefix(module.Runtime, "kubernetes:") {
		namespace := strings.TrimPrefix(module.Runtime, "kubernetes:")
		client := kubernetes.Client{Executable: os.Getenv("BITPACK_KUBECTL"), Namespace: namespace}
		return client.Logs(context.Background(), module.Name, output)
	}
	if module.ProcessID != "" {
		return processManager().Logs(module.ProcessID, output)
	}
	path := dataPath("logs", "modules", module.ID+".log")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		fmt.Fprintf(output, "No logs found for %s.\n", module.ID)
		return nil
	}
	if err != nil {
		return fmt.Errorf("read module logs: %w", err)
	}
	_, err = output.Write(data)
	return err
}

func bitpackSnapshot(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: snapshot <module-id-or-name>")
	}
	store := moduleStore()
	module, err := store.Find(args[0])
	if err != nil {
		return err
	}
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	path := dataPath("snapshots", module.ID, id+".json")
	record := snapshotRecord{ID: id, Module: module.ID, Created: time.Now().Unix(), State: *module}
	if err := writeJSON(path, record, 0o600); err != nil {
		return fmt.Errorf("write module snapshot: %w", err)
	}
	fmt.Fprintf(output, "Snapshot created: %s\n", id)
	return nil
}

func bitpackRestore(output io.Writer, args []string) error {
	if len(args) != 2 || !validResourceName.MatchString(args[1]) {
		return fmt.Errorf("usage: restore <module-id-or-name> <snapshot-id>")
	}
	store := moduleStore()
	module, err := store.Find(args[0])
	if err != nil {
		return err
	}
	path := dataPath("snapshots", module.ID, args[1]+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read snapshot %s: %w", args[1], err)
	}
	var snapshot snapshotRecord
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("decode snapshot %s: %w", args[1], err)
	}
	if snapshot.Module != module.ID || snapshot.State.ID != module.ID {
		return fmt.Errorf("snapshot %s does not belong to module %s", args[1], module.ID)
	}
	snapshot.State.Status = "running"
	if err := store.Update(&snapshot.State); err != nil {
		return err
	}
	fmt.Fprintf(output, "Module state restored: %s\n", module.ID)
	return nil
}
