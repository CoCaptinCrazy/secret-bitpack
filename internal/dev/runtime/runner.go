package runtime

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/dev/parser"
	"github.com/CoCaptinCrazy/secret-bitpack/internal/runtime/supervisor"
	bitpacksecrets "github.com/CoCaptinCrazy/secret-bitpack/internal/security/secrets"
)

func Load(path string) (*parser.Environment, string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("resolve development file: %w", err)
	}
	file, err := os.Open(absolutePath)
	if err != nil {
		return nil, "", fmt.Errorf("open development file %s: %w", path, err)
	}
	defer file.Close()

	environment, err := parser.Parse(file)
	if err != nil {
		return nil, "", fmt.Errorf("parse development file %s: %w", path, err)
	}
	workspace := environment.Workspace
	if !filepath.IsAbs(workspace) {
		workspace = filepath.Join(filepath.Dir(absolutePath), workspace)
	}
	workspace, err = filepath.Abs(workspace)
	if err != nil {
		return nil, "", fmt.Errorf("resolve development workspace: %w", err)
	}
	info, err := os.Stat(workspace)
	if err != nil {
		return nil, "", fmt.Errorf("open development workspace %s: %w", workspace, err)
	}
	if !info.IsDir() {
		return nil, "", fmt.Errorf("development workspace %s is not a directory", workspace)
	}
	return environment, workspace, nil
}

func Run(
	ctx context.Context,
	environment *parser.Environment,
	workspace string,
	manager *supervisor.Manager,
	input io.Reader,
	output, errorOutput io.Writer,
) (returnErr error) {
	if manager == nil {
		return fmt.Errorf("process supervisor is required")
	}
	for _, service := range environment.Services {
		if service.Command == "" {
			return fmt.Errorf("service %q requires a host command (service %s = <command>)",
				service.Name, service.Name)
		}
	}
	variables, err := resolvedEnvironment(environment.Variables, environment.Secrets)
	if err != nil {
		return err
	}
	serviceIDs := make([]string, 0, len(environment.Services))
	appIDs := make([]string, 0, len(environment.Shell))
	startedIDs := make([]string, 0, len(environment.Services)+len(environment.Shell))
	defer func() {
		for index := len(startedIDs) - 1; index >= 0; index-- {
			if err := manager.Stop(startedIDs[index]); err != nil {
				if returnErr == nil {
					returnErr = fmt.Errorf("stop development process %s: %w", startedIDs[index], err)
				} else {
					returnErr = fmt.Errorf("%w; stop development process %s: %v",
						returnErr, startedIDs[index], err)
				}
			}
		}
	}()

	for index, service := range environment.Services {
		id := processID(environment.Name, "service", index, service.Name)
		record, err := startProcess(manager, id, service.Name, service.Command, workspace,
			environment.Variables, environment.Secrets)
		if err != nil {
			return fmt.Errorf("start development service %q: %w", service.Name, err)
		}
		serviceIDs = append(serviceIDs, id)
		startedIDs = append(startedIDs, id)
		fmt.Fprintf(output, "Service started: %s (pid %d)\n", service.Name, record.PID)
	}
	if err := RunCommands(ctx, environment.Startup, variables, workspace, input, output, errorOutput); err != nil {
		return fmt.Errorf("development startup failed: %w", err)
	}
	for index, command := range environment.Shell {
		id := processID(environment.Name, "app", index, "shell")
		record, err := startProcess(manager, id, fmt.Sprintf("%s shell %d", environment.Name, index+1),
			command, workspace, environment.Variables, environment.Secrets)
		if err != nil {
			return fmt.Errorf("start development command: %w", err)
		}
		appIDs = append(appIDs, id)
		startedIDs = append(startedIDs, id)
		fmt.Fprintf(output, "Development process started: %s (pid %d)\n", id, record.PID)
	}
	if len(startedIDs) == 0 {
		return nil
	}
	if len(environment.Watch) == 0 {
		fmt.Fprintln(output, "Development processes are supervised; press Ctrl+C to stop.")
	}
	previous, err := Snapshot(workspace, environment.Watch)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	pendingWatchRestart := false
	lastWatchChange := time.Time{}
	lastRestart := make(map[string]time.Time)
	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(output, "Stopping development processes.")
			return nil
		case now := <-ticker.C:
			current, err := Snapshot(workspace, environment.Watch)
			if err != nil {
				return err
			}
			changed := ChangedPaths(previous, current)
			if len(changed) > 0 {
				sort.Strings(changed)
				fmt.Fprintf(output, "[WATCH] Changed: %s\n", strings.Join(changed, ", "))
				previous = current
				pendingWatchRestart = len(appIDs) > 0
				lastWatchChange = now
			} else if pendingWatchRestart && now.Sub(lastWatchChange) >= 450*time.Millisecond {
				for _, id := range appIDs {
					record, err := manager.Restart(id)
					if err != nil {
						return fmt.Errorf("restart development process %s after file change: %w", id, err)
					}
					fmt.Fprintf(output, "[RESTART] %s (pid %d)\n", id, record.PID)
				}
				pendingWatchRestart = false
			}
			for _, id := range append(append([]string(nil), serviceIDs...), appIDs...) {
				running, err := manager.IsRunning(id)
				if err != nil {
					return err
				}
				if running || now.Sub(lastRestart[id]) < time.Second {
					continue
				}
				record, err := manager.Restart(id)
				if err != nil {
					return fmt.Errorf("restart exited development process %s: %w", id, err)
				}
				lastRestart[id] = now
				fmt.Fprintf(output, "[SUPERVISE] Restarted %s (pid %d)\n", id, record.PID)
			}
		}
	}
}

func startProcess(
	manager *supervisor.Manager,
	id, name, line, workspace string,
	variables map[string]string,
	secrets []string,
) (*supervisor.Record, error) {
	executable, arguments, err := supervisor.ShellCommand(line)
	if err != nil {
		return nil, err
	}
	if _, err := manager.Get(id); err == nil {
		running, checkErr := manager.IsRunning(id)
		if checkErr != nil {
			return nil, checkErr
		}
		if running {
			if err := manager.Stop(id); err != nil {
				return nil, err
			}
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return manager.Start(id, supervisor.Spec{
		Name: name, Executable: executable, Arguments: arguments,
		WorkingDir: workspace, Environment: variables, Secrets: secrets,
	})
}

func processID(environment, kind string, index int, name string) string {
	prefix := supervisor.IDPart(environment)
	if prefix == "" {
		prefix = "development"
	}
	service := supervisor.IDPart(name)
	if service == "" {
		service = "process"
	}
	return supervisor.IDWithIndex(prefix+"-"+kind+"-"+service, index+1)
}

func resolvedEnvironment(variables map[string]string, secretNames []string) (map[string]string, error) {
	result := make(map[string]string, len(variables)+len(secretNames))
	for key, value := range variables {
		result[key] = value
	}
	if len(secretNames) > 0 {
		values, err := bitpacksecrets.ResolveFromEnvironment(secretNames)
		if err != nil {
			return nil, err
		}
		for name, value := range values {
			key := bitpacksecrets.EnvironmentKey(name)
			if _, exists := result[key]; exists {
				return nil, fmt.Errorf("secret %q conflicts with environment variable %q", name, key)
			}
			result[key] = value
		}
	}
	return result, nil
}

func RunCommands(ctx context.Context, commands []string, variables map[string]string, workspace string, input io.Reader, output, errorOutput io.Writer) error {
	env := commandEnvironment(variables)
	for _, line := range commands {
		if err := execute(ctx, workspace, line, env, input, output, errorOutput); err != nil {
			return fmt.Errorf("command %q failed: %w", line, err)
		}
	}
	return nil
}

func execute(ctx context.Context, workspace, line string, environment []string, input io.Reader, output, errorOutput io.Writer) error {
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		executable, err := exec.LookPath("cmd.exe")
		if err != nil {
			return fmt.Errorf("find cmd.exe: %w", err)
		}
		command = exec.CommandContext(ctx, executable, "/D", "/S", "/C", line)
	} else {
		executable, err := exec.LookPath("sh")
		if err != nil {
			return fmt.Errorf("find sh: %w", err)
		}
		command = exec.CommandContext(ctx, executable, "-c", line)
	}
	command.Dir = workspace
	command.Env = environment
	command.Stdin = input
	command.Stdout = output
	command.Stderr = errorOutput
	return command.Run()
}

func commandEnvironment(overrides map[string]string) []string {
	values := make(map[string]string)
	for _, value := range os.Environ() {
		key, contents, ok := strings.Cut(value, "=")
		if ok {
			values[key] = contents
		}
	}
	for key, value := range overrides {
		values[key] = value
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, key+"="+values[key])
	}
	return result
}
