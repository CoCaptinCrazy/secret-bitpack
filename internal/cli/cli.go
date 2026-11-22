package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"secret-bitpack/internal/core"
	"secret-bitpack/internal/core/bootstrap"
	devruntime "secret-bitpack/internal/dev/runtime"
	"secret-bitpack/internal/extensions"
	"secret-bitpack/internal/orchestration/cluster"
)

func Execute(args []string, output io.Writer) error {
	return ExecuteWithIO(args, os.Stdin, output)
}

func ExecuteWithIO(args []string, input io.Reader, output io.Writer) error {
	bitpack, err := bootstrap.Initialize(output)
	if err != nil {
		return err
	}
	if err := setDataDirectory(bitpack.Config.DataDir); err != nil {
		return err
	}
	return execute(bitpack, args, input, output)
}

func execute(bitpack *core.Bitpack, args []string, input io.Reader, output io.Writer) error {
	if len(args) == 0 {
		args = []string{"help"}
	}

	command := strings.ToLower(args[0])
	commandArgs := args[1:]
	switch command {
	case "help", "-h", "--help":
		printHelp(output)
	case "version", "--version", "-v":
		fmt.Fprintf(output, "Bitpack v%s\n", bitpack.Version)
	case "config":
		return bitpackConfig(bitpack, output, commandArgs)
	case "forge":
		return bitpackForge(output, commandArgs)
	case "pack":
		return bitpackPackage(output, commandArgs)
	case "unpack":
		return bitpackUnpack(output, commandArgs)
	case "image":
		return bitpackImage(output, commandArgs)
	case "verify":
		return bitpackVerify(output, commandArgs)
	case "blueprints":
		return bitpackBlueprints(output)
	case "inspect":
		return bitpackInspect(output, commandArgs)
	case "export":
		return bitpackExport(output, commandArgs)
	case "import":
		return bitpackImport(output, commandArgs)
	case "ignite":
		return bitpackIgnite(output, commandArgs)
	case "modules":
		return bitpackModules(output)
	case "stop":
		return bitpackStop(output, commandArgs)
	case "restart":
		return bitpackRestart(output, commandArgs)
	case "rm", "remove":
		return bitpackRemove(output, commandArgs)
	case "inspect-module":
		return bitpackInspectModule(output, commandArgs)
	case "logs":
		return bitpackLogs(output, commandArgs)
	case "snapshot":
		return bitpackSnapshot(output, commandArgs)
	case "restore":
		return bitpackRestore(output, commandArgs)
	case "runtime":
		return bitpackRuntime(output)
	case "health":
		return bitpackHealth(output, commandArgs)
	case "backup":
		return bitpackBackup(output)
	case "metrics":
		return bitpackMetrics(output)
	case "profile":
		return bitpackProfile(output)
	case "alert":
		return bitpackAlert(output, 95)
	case "alerts":
		return bitpackAlert(output, 100)
	case "events":
		return bitpackEvents(output)
	case "secrets":
		return bitpackSecrets(output, commandArgs)
	case "mesh":
		return bitpackMesh(output, commandArgs)
	case "meshes":
		return bitpackMeshes(output)
	case "bridges":
		return bitpackBridges(output)
	case "route":
		return bitpackRoutes(output)
	case "nodes":
		return bitpackNodes(output)
	case "policy":
		return bitpackPolicy(output, commandArgs)
	case "search":
		return bitpackSearch(output, commandArgs)
	case "publish":
		return bitpackPublish(output, commandArgs)
	case "acquire":
		return bitpackAcquire(output, commandArgs)
	case "tags":
		return bitpackTags(output, commandArgs)
	case "login":
		return bitpackLogin(bitpack, output, commandArgs)
	case "logout":
		return bitpackLogout(bitpack, output)
	case "assemble":
		return bitpackAssemble(output, commandArgs)
	case "assemblies":
		return bitpackAssemblies(output)
	case "destroy":
		return bitpackDestroy(output, commandArgs)
	case "scale":
		return bitpackScale(output, commandArgs)
	case "rollout":
		return bitpackRollout(output, commandArgs, false)
	case "rollback":
		return bitpackRollout(output, commandArgs, true)
	case "watch":
		return bitpackWatch(output, commandArgs)
	case "batch":
		return bitpackBatch(output, commandArgs)
	case "pipeline":
		return bitpackBatch(output, commandArgs)
	case "batch-status":
		return bitpackBatchStatus(output)
	case "shell":
		return runShell(input, output, func(line []string) error {
			return execute(bitpack, line, input, output)
		})
	case "plugin":
		return runPlugin(commandArgs, input, output)
	case "cluster":
		return bitpackCluster(output, commandArgs, bitpack.Version)
	case "clusters":
		return bitpackCluster(output, append([]string{"list"}, commandArgs...), bitpack.Version)
	case "dev":
		return bitpackDev(input, output, commandArgs)
	default:
		return fmt.Errorf("unknown command %q (run 'help' for available commands)", command)
	}
	return nil
}

func printHelp(output io.Writer) {
	fmt.Fprintln(output, `Bitpack - Polyglot project and deployment tool

Usage:
  bitpack <command> [arguments]

Commands:
  help, version, config [get|set]
  forge [Bitfile], pack [Bitfile] [output.bpkg], unpack <file.bpkg> [destination]
  image build [Bitfile] --tag <registry/repository:tag>, image push <tag>
  verify [blueprint|Bitfile], blueprints, inspect, export, import
  ignite <blueprint|Bitfile>, modules, stop, restart, rm, inspect-module, logs
  snapshot, restore
  assemble <blueprint|file.bpack>, assemblies, destroy, scale, rollout, rollback
  publish, acquire, search, tags, login, logout
  dev [inspect|run] [project.bpdev], shell, watch
  batch, batch-status, pipeline
  mesh, meshes, bridges, route
  health, policy, alert, alerts, metrics, profile, events, secrets
  cluster, clusters, nodes, runtime, backup
  plugin run <name>       Run an explicitly selected Python extension

Use 'exit' or 'quit' to leave the shell.
Project commands use the configured operating-system shell and installed toolchains.
Python is only started for an explicitly selected plugin.`)
}

func bitpackCluster(output io.Writer, args []string, version string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: cluster <create|list|inspect|delete> [name]")
	}
	store := cluster.NewStore(dataPath("clusters"), version)
	switch args[0] {
	case "create":
		if len(args) != 2 {
			return fmt.Errorf("usage: cluster create <name>")
		}
		record, err := store.Create(args[1])
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "Cluster definition created: %s (pending; no infrastructure provisioned)\n",
			record.Metadata.Name)
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("usage: cluster list")
		}
		records, err := store.List()
		if err != nil {
			return err
		}
		if len(records) == 0 {
			fmt.Fprintln(output, "No cluster definitions.")
			return nil
		}
		for _, record := range records {
			fmt.Fprintf(output, "%s\t%s\n", record.Metadata.Name, record.Status.Phase)
		}
	case "inspect":
		if len(args) != 2 {
			return fmt.Errorf("usage: cluster inspect <name>")
		}
		record, err := store.Get(args[1])
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "Name: %s\nStatus: %s\nHealthy: %t\nNodes: %d\n",
			record.Metadata.Name, record.Status.Phase, record.Status.Healthy, len(record.Nodes))
	case "delete", "rm":
		if len(args) != 2 {
			return fmt.Errorf("usage: cluster delete <name>")
		}
		if err := store.Delete(args[1]); err != nil {
			return err
		}
		fmt.Fprintf(output, "Cluster definition deleted: %s\n", args[1])
	default:
		return fmt.Errorf("unknown cluster action %q; use create, list, inspect, or delete", args[0])
	}
	return nil
}

func bitpackDev(input io.Reader, output io.Writer, args []string) error {
	action := "inspect"
	path := "project.bpdev"
	if len(args) > 0 {
		switch args[0] {
		case "inspect", "run":
			action = args[0]
			if len(args) > 1 {
				path = args[1]
			}
			if len(args) > 2 {
				return fmt.Errorf("usage: dev [inspect|run] [project.bpdev]")
			}
		default:
			if len(args) > 1 {
				return fmt.Errorf("usage: dev [inspect|run] [project.bpdev]")
			}
			path = args[0]
		}
	}
	environment, workspace, err := devruntime.Load(path)
	if err != nil {
		return err
	}
	fmt.Fprintf(output, "Development environment: %s v%s\nRuntime: %s\nWorkspace: %s\n",
		environment.Name, environment.Version, environment.Runtime, workspace)
	if action == "run" {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		return devruntime.Run(ctx, environment, workspace, processManager(), input, output, os.Stderr)
	}
	for _, service := range environment.Services {
		fmt.Fprintf(output, "Service: %s\n", service.Name)
		if service.Command != "" {
			fmt.Fprintf(output, "  Command: %s\n", service.Command)
		}
	}
	for _, port := range environment.Ports {
		fmt.Fprintf(output, "Port: %d\n", port)
	}
	for _, path := range environment.Watch {
		fmt.Fprintf(output, "Watch: %s\n", path)
	}
	for _, command := range environment.Startup {
		fmt.Fprintf(output, "Startup: %s\n", command)
	}
	for _, command := range environment.Shell {
		fmt.Fprintf(output, "Shell: %s\n", command)
	}
	return nil
}

func runPlugin(args []string, input io.Reader, output io.Writer) error {
	if len(args) < 2 || args[0] != "run" {
		return fmt.Errorf("usage: plugin run <name> [arguments...]")
	}
	return extensions.Run(input, output, os.Stderr, "plugins", args[1], args[2:]...)
}
