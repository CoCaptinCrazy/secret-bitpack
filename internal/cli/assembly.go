package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	assemblyparser "secret-bitpack/internal/assembly/parser"
	devruntime "secret-bitpack/internal/dev/runtime"
	"secret-bitpack/internal/orchestration/kubernetes"
	"secret-bitpack/internal/runtime/modules"
	bitpacksecrets "secret-bitpack/internal/security/secrets"
)

type assemblyDefinition struct {
	Name      string          `json:"name"`
	Source    string          `json:"source"`
	Created   int64           `json:"created"`
	Provider  string          `json:"provider,omitempty"`
	Namespace string          `json:"namespace,omitempty"`
	Status    string          `json:"status,omitempty"`
	Resources []assemblyEntry `json:"resources"`
}

type assemblyEntry struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type deploymentRevision struct {
	Name       string `json:"name"`
	Generation int    `json:"generation"`
	Updated    int64  `json:"updated"`
}

func bitpackAssemble(output io.Writer, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: assemble <blueprint-or-file.bpack> [--dry-run] [--namespace name]")
	}
	if filepath.Ext(args[0]) == ".bpack" {
		return bitpackAssembleDefinition(output, args)
	}
	if len(args) != 1 {
		return fmt.Errorf("usage: assemble <blueprint-or-file.bpack>")
	}
	parsed, workspace, _, err := loadBlueprintInput(args[0])
	if err != nil {
		return err
	}
	if len(parsed.Assemble) == 0 {
		fmt.Fprintf(output, "No assemble commands defined for %s.\n", parsed.Name)
		return nil
	}
	environment, err := blueprintEnvironment(parsed)
	if err != nil {
		return err
	}
	if err := devruntime.RunCommands(
		context.Background(), parsed.Assemble, environment, workspace, os.Stdin, output, os.Stderr,
	); err != nil {
		return fmt.Errorf("assemble blueprint %q: %w", parsed.Name, err)
	}
	fmt.Fprintf(output, "Assemble commands completed: %s\n", parsed.Name)
	return nil
}

func bitpackAssembleDefinition(output io.Writer, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: assemble <file.bpack> [--dry-run] [--namespace name]")
	}
	path := args[0]
	namespace := os.Getenv("BITPACK_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}
	dryRun := false
	for index := 1; index < len(args); index++ {
		switch args[index] {
		case "--dry-run":
			dryRun = true
		case "--namespace":
			if index+1 == len(args) {
				return fmt.Errorf("--namespace requires a value")
			}
			index++
			namespace = args[index]
		default:
			return fmt.Errorf("unknown assemble option %q", args[index])
		}
	}
	if !validKubernetesName(namespace) {
		return fmt.Errorf("invalid Kubernetes namespace %q", namespace)
	}
	client := kubernetes.Client{
		Executable: os.Getenv("BITPACK_KUBECTL"),
		Namespace:  namespace,
	}
	return deployBPack(context.Background(), path, namespace, dryRun, client, output)
}

func deployBPack(ctx context.Context, path, namespace string, dryRun bool, client kubernetes.Client, output io.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open BPack %s: %w", path, err)
	}
	assembly, parseErr := assemblyparser.ParseBPack(file)
	closeErr := file.Close()
	if parseErr != nil {
		return fmt.Errorf("parse BPack %s: %w", path, parseErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close BPack %s: %w", path, closeErr)
	}
	secretNames := make([]string, 0)
	for _, module := range assembly.Modules {
		secretNames = append(secretNames, module.Secrets...)
	}
	secretValues := map[string]string{}
	if len(secretNames) > 0 {
		secretValues, err = bitpacksecrets.ResolveFromEnvironment(secretNames)
		if err != nil {
			return err
		}
	}
	objects, err := kubernetes.Compile(assembly, secretValues)
	if err != nil {
		return fmt.Errorf("compile BPack %s: %w", path, err)
	}
	if dryRun {
		return printKubernetesObjects(output, objects, secretValues)
	}
	if err := client.Apply(ctx, objects, secretValues, output); err != nil {
		return err
	}
	root := dataPath("assemblies", assembly.Name)
	if err := os.MkdirAll(root, 0o700); err != nil {
		return fmt.Errorf("create assembly state directory after Kubernetes apply: %w", err)
	}
	definition := assemblyDefinition{
		Name: assembly.Name, Source: path, Created: time.Now().Unix(),
		Provider: "kubernetes", Namespace: namespace, Status: "active",
		Resources: make([]assemblyEntry, 0, len(assembly.Modules)),
	}
	for _, module := range assembly.Modules {
		definition.Resources = append(definition.Resources, assemblyEntry{Kind: "module", Name: module.Name})
		var autoscale *assemblyparser.Autoscale
		for index := range assembly.Autoscale {
			if assembly.Autoscale[index].Target == module.Name {
				autoscale = &assembly.Autoscale[index]
				break
			}
		}
		if err := upsertKubernetesModule(module.Name, assembly.Name, namespace, module.Replicas, autoscale); err != nil {
			return fmt.Errorf("record Kubernetes module %q after apply: %w", module.Name, err)
		}
	}
	if err := writeJSON(filepath.Join(root, "assembly.json"), definition, 0o600); err != nil {
		return fmt.Errorf("record applied Kubernetes assembly: %w", err)
	}
	fmt.Fprintf(output, "Kubernetes assembly applied: %s in namespace %s\n", assembly.Name, namespace)
	return nil
}

func printKubernetesObjects(output io.Writer, objects []map[string]any, secretValues map[string]string) error {
	for _, object := range objects {
		data, err := json.MarshalIndent(object, "", "  ")
		if err != nil {
			return fmt.Errorf("encode Kubernetes dry-run resource: %w", err)
		}
		rendered := string(data)
		for _, value := range secretValues {
			if value != "" {
				rendered = strings.ReplaceAll(rendered, value, "[REDACTED]")
			}
		}
		if _, err := fmt.Fprintln(output, rendered); err != nil {
			return fmt.Errorf("write Kubernetes dry-run resource: %w", err)
		}
	}
	return nil
}

func validKubernetesName(value string) bool {
	return value != "" && validKubernetesResourceName.MatchString(value)
}

var validKubernetesResourceName = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

func upsertKubernetesModule(name, assemblyName, namespace string, replicas int, autoscale *assemblyparser.Autoscale) error {
	store := moduleStore()
	existingModules, err := store.List()
	if err != nil {
		return err
	}
	setAutoscale := func(module *modules.Module) {
		module.AutoscaleMin = 0
		module.AutoscaleMax = 0
		module.AutoscaleCPU = 0
		if autoscale != nil {
			module.AutoscaleMin = autoscale.Min
			module.AutoscaleMax = autoscale.Max
			module.AutoscaleCPU = autoscale.CPU
		}
	}
	for index := range existingModules {
		if existingModules[index].Name != name {
			continue
		}
		existingModules[index].Blueprint = assemblyName
		existingModules[index].Runtime = "kubernetes:" + namespace
		existingModules[index].Status = "running"
		existingModules[index].DesiredReplicas = replicas
		existingModules[index].Started = time.Now().Unix()
		setAutoscale(&existingModules[index])
		return store.Update(&existingModules[index])
	}
	module, err := store.Create(name, assemblyName, "kubernetes:"+namespace)
	if err != nil {
		return err
	}
	module.DesiredReplicas = replicas
	setAutoscale(module)
	return store.Update(module)
}

func bitpackAssemblies(output io.Writer) error {
	return listDirectories(output, dataPath("assemblies"), "No assemblies found.")
}

func bitpackDestroy(output io.Writer, args []string) error {
	if len(args) != 1 || !validResourceName.MatchString(args[0]) {
		return fmt.Errorf("usage: destroy <assembly>")
	}
	path := dataPath("assemblies", args[0])
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("assembly %q not found", args[0])
	} else if err != nil {
		return fmt.Errorf("inspect assembly %q: %w", args[0], err)
	}
	definitionPath := filepath.Join(path, "assembly.json")
	if data, err := os.ReadFile(definitionPath); err == nil {
		var definition assemblyDefinition
		if err := json.Unmarshal(data, &definition); err != nil {
			return fmt.Errorf("decode assembly state: %w", err)
		}
		if definition.Provider == "kubernetes" {
			client := kubernetes.Client{
				Executable: os.Getenv("BITPACK_KUBECTL"),
				Namespace:  definition.Namespace,
			}
			if err := client.DeleteAssembly(context.Background(), definition.Name, output); err != nil {
				return err
			}
			modules, err := moduleStore().List()
			if err != nil {
				return fmt.Errorf("list modules after Kubernetes deletion: %w", err)
			}
			for index := range modules {
				if modules[index].Blueprint == definition.Name {
					if err := moduleStore().Delete(&modules[index]); err != nil {
						return fmt.Errorf("remove module record %q: %w", modules[index].ID, err)
					}
				}
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read assembly state: %w", err)
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("destroy assembly definition %s: %w", args[0], err)
	}
	fmt.Fprintf(output, "Assembly destroyed: %s\n", args[0])
	return nil
}

func bitpackScale(output io.Writer, args []string) error {
	if len(args) != 2 || !validResourceName.MatchString(args[0]) {
		return fmt.Errorf("usage: scale <service> <replicas>")
	}
	replicas, err := strconv.Atoi(args[1])
	if err != nil || replicas < 0 {
		return fmt.Errorf("replicas must be a non-negative integer")
	}
	module, namespace, err := findKubernetesModule(args[0])
	if err != nil {
		return err
	}
	if module != nil {
		client := kubernetes.Client{Executable: os.Getenv("BITPACK_KUBECTL"), Namespace: namespace}
		if err := client.Scale(context.Background(), module.Name, replicas, output); err != nil {
			return err
		}
		module.DesiredReplicas = replicas
		if replicas == 0 {
			module.Status = "stopped"
		} else {
			module.Status = "running"
		}
		return moduleStore().Update(module)
	}
	path := dataPath("scaling", args[0]+".json")
	if err := writeJSON(path, map[string]any{
		"service": args[0], "replicas": replicas, "updated": time.Now().Unix(),
	}, 0o600); err != nil {
		return fmt.Errorf("save scaling configuration: %w", err)
	}
	fmt.Fprintf(output, "Scale recorded: %s -> %d replicas (not applied to infrastructure)\n",
		args[0], replicas)
	return nil
}

func bitpackRollout(output io.Writer, args []string, rollback bool) error {
	if len(args) != 1 || !validResourceName.MatchString(args[0]) {
		action := "rollout"
		if rollback {
			action = "rollback"
		}
		return fmt.Errorf("usage: %s <assembly>", action)
	}
	assemblyPath := dataPath("assemblies", args[0], "assembly.json")
	if data, err := os.ReadFile(assemblyPath); err == nil {
		var definition assemblyDefinition
		if err := json.Unmarshal(data, &definition); err != nil {
			return fmt.Errorf("decode assembly state: %w", err)
		}
		if definition.Provider == "kubernetes" {
			client := kubernetes.Client{
				Executable: os.Getenv("BITPACK_KUBECTL"),
				Namespace:  definition.Namespace,
			}
			for _, resource := range definition.Resources {
				if resource.Kind != "module" {
					continue
				}
				var err error
				if rollback {
					err = client.Rollback(context.Background(), resource.Name, output)
				} else {
					err = client.Restart(context.Background(), resource.Name, output)
				}
				if err != nil {
					return err
				}
			}
			action := "Kubernetes rollout restarted"
			if rollback {
				action = "Kubernetes rollout rolled back"
			}
			fmt.Fprintf(output, "%s: %s\n", action, definition.Name)
			return nil
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read assembly state: %w", err)
	}
	path := dataPath("deployments", args[0]+".json")
	revision := deploymentRevision{Name: args[0], Generation: 0}
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &revision); err != nil {
			return fmt.Errorf("decode deployment revision: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read deployment revision: %w", err)
	}
	if rollback {
		if revision.Generation == 0 {
			return fmt.Errorf("no rollout exists for assembly %q", args[0])
		}
		revision.Generation--
	} else {
		revision.Generation++
	}
	revision.Updated = time.Now().Unix()
	if err := writeJSON(path, revision, 0o600); err != nil {
		return fmt.Errorf("save deployment revision: %w", err)
	}
	action := "Rollout"
	if rollback {
		action = "Rollback"
	}
	fmt.Fprintf(output, "%s revision recorded: %s generation %d (not applied to infrastructure)\n",
		action, args[0], revision.Generation)
	return nil
}

func findKubernetesModule(name string) (*modules.Module, string, error) {
	list, err := moduleStore().List()
	if err != nil {
		return nil, "", err
	}
	for index := range list {
		module := &list[index]
		if module.Name != name || !strings.HasPrefix(module.Runtime, "kubernetes:") {
			continue
		}
		namespace := strings.TrimPrefix(module.Runtime, "kubernetes:")
		if namespace == "" {
			return nil, "", fmt.Errorf("Kubernetes module %q has no namespace", name)
		}
		return module, namespace, nil
	}
	return nil, "", nil
}

func readKubernetesAssembly(name string) (*assemblyDefinition, bool, error) {
	if !validResourceName.MatchString(name) {
		return nil, false, nil
	}
	path := dataPath("assemblies", name, "assembly.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read assembly state: %w", err)
	}
	var definition assemblyDefinition
	if err := json.Unmarshal(data, &definition); err != nil {
		return nil, false, fmt.Errorf("decode assembly state: %w", err)
	}
	return &definition, definition.Provider == "kubernetes", nil
}

func stopKubernetesAssembly(output io.Writer, definition *assemblyDefinition) error {
	client := kubernetes.Client{
		Executable: os.Getenv("BITPACK_KUBECTL"),
		Namespace:  definition.Namespace,
	}
	for _, resource := range definition.Resources {
		if resource.Kind != "module" {
			continue
		}
		module, err := moduleStore().Find(resource.Name)
		if err != nil {
			return err
		}
		if err := client.Stop(context.Background(), module.Name, output); err != nil {
			return err
		}
		module.Status = "stopped"
		if err := moduleStore().Update(module); err != nil {
			return err
		}
	}
	definition.Status = "stopped"
	if err := writeJSON(dataPath("assemblies", definition.Name, "assembly.json"), definition, 0o600); err != nil {
		return fmt.Errorf("record stopped assembly: %w", err)
	}
	fmt.Fprintf(output, "Kubernetes assembly stopped: %s\n", definition.Name)
	return nil
}

func restartKubernetesAssembly(output io.Writer, definition *assemblyDefinition) error {
	client := kubernetes.Client{
		Executable: os.Getenv("BITPACK_KUBECTL"),
		Namespace:  definition.Namespace,
	}
	for _, resource := range definition.Resources {
		if resource.Kind != "module" {
			continue
		}
		module, err := moduleStore().Find(resource.Name)
		if err != nil {
			return err
		}
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
			if err := client.ApplyAutoscaler(context.Background(), definition.Name, module.Name,
				module.AutoscaleMin, module.AutoscaleMax, module.AutoscaleCPU, output); err != nil {
				return err
			}
		}
		module.Status = "running"
		module.DesiredReplicas = replicas
		module.Started = time.Now().Unix()
		if err := moduleStore().Update(module); err != nil {
			return err
		}
	}
	definition.Status = "active"
	if err := writeJSON(dataPath("assemblies", definition.Name, "assembly.json"), definition, 0o600); err != nil {
		return fmt.Errorf("record restarted assembly: %w", err)
	}
	fmt.Fprintf(output, "Kubernetes assembly restarted: %s\n", definition.Name)
	return nil
}
