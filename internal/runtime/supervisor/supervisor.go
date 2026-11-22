package supervisor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	bitpacksecrets "secret-bitpack/internal/security/secrets"
)

var validID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`)

type Spec struct {
	Name        string            `json:"name"`
	Executable  string            `json:"executable"`
	Arguments   []string          `json:"arguments,omitempty"`
	WorkingDir  string            `json:"working_dir"`
	Environment map[string]string `json:"environment,omitempty"`
	Secrets     []string          `json:"secrets,omitempty"`
}

type Record struct {
	ID      string    `json:"id"`
	PID     int       `json:"pid"`
	Status  string    `json:"status"`
	Started time.Time `json:"started"`
	Spec    Spec      `json:"spec"`
	LogPath string    `json:"log_path"`
}

type Manager struct {
	Directory    string
	LogDirectory string
}

func New(directory, logDirectory string) *Manager {
	return &Manager{Directory: directory, LogDirectory: logDirectory}
}

func (m *Manager) Start(id string, spec Spec) (*Record, error) {
	if !validID.MatchString(id) {
		return nil, fmt.Errorf("invalid supervised process ID %q", id)
	}
	if spec.Name == "" || spec.Executable == "" || spec.WorkingDir == "" {
		return nil, fmt.Errorf("process name, executable, and working directory are required")
	}
	if info, err := os.Stat(spec.WorkingDir); err != nil {
		return nil, fmt.Errorf("inspect process working directory: %w", err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("process working directory %s is not a directory", spec.WorkingDir)
	}
	if err := os.MkdirAll(m.Directory, 0o700); err != nil {
		return nil, fmt.Errorf("create process state directory: %w", err)
	}
	if err := os.MkdirAll(m.LogDirectory, 0o700); err != nil {
		return nil, fmt.Errorf("create process log directory: %w", err)
	}
	if existing, err := m.Get(id); err == nil {
		running, checkErr := processRunning(existing.PID)
		if checkErr != nil {
			return nil, checkErr
		}
		if running {
			return nil, fmt.Errorf("process %q is already running", id)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	environment, err := processEnvironment(spec)
	if err != nil {
		return nil, err
	}
	logPath := filepath.Join(m.LogDirectory, id+".log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open process log: %w", err)
	}
	command := exec.Command(spec.Executable, spec.Arguments...)
	command.Dir = spec.WorkingDir
	command.Env = environment
	command.Stdout = logFile
	command.Stderr = logFile
	configureProcess(command)
	if err := command.Start(); err != nil {
		if closeErr := logFile.Close(); closeErr != nil {
			return nil, fmt.Errorf("start process: %w; close process log: %v", err, closeErr)
		}
		return nil, fmt.Errorf("start process %q: %w", spec.Name, err)
	}
	pid := command.Process.Pid
	if err := command.Process.Release(); err != nil {
		_ = terminateProcess(pid)
		_ = logFile.Close()
		return nil, fmt.Errorf("release process handle for %q: %w", spec.Name, err)
	}
	if err := logFile.Close(); err != nil {
		_ = terminateProcess(pid)
		return nil, fmt.Errorf("close process log: %w", err)
	}
	record := &Record{
		ID: id, PID: pid, Status: "running", Started: time.Now().UTC(),
		Spec: spec, LogPath: logPath,
	}
	if err := m.save(record); err != nil {
		_ = terminateProcess(pid)
		return nil, err
	}
	return record, nil
}

func (m *Manager) Stop(id string) error {
	record, err := m.Get(id)
	if err != nil {
		return err
	}
	running, err := processRunning(record.PID)
	if err != nil {
		return err
	}
	if running {
		if err := terminateProcess(record.PID); err != nil {
			return fmt.Errorf("stop process %q: %w", record.Spec.Name, err)
		}
	}
	record.Status = "stopped"
	if err := m.save(record); err != nil {
		return fmt.Errorf("record stopped process: %w", err)
	}
	return nil
}

func (m *Manager) Restart(id string) (*Record, error) {
	record, err := m.Get(id)
	if err != nil {
		return nil, err
	}
	if err := m.Stop(id); err != nil {
		return nil, err
	}
	return m.Start(id, record.Spec)
}

func (m *Manager) Get(id string) (*Record, error) {
	if !validID.MatchString(id) {
		return nil, fmt.Errorf("invalid supervised process ID %q", id)
	}
	data, err := os.ReadFile(m.statePath(id))
	if err != nil {
		return nil, err
	}
	var record Record
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("decode process state %s: %w", id, err)
	}
	return &record, nil
}

func (m *Manager) IsRunning(id string) (bool, error) {
	record, err := m.Get(id)
	if err != nil {
		return false, err
	}
	running, err := processRunning(record.PID)
	if err != nil {
		return false, err
	}
	if !running && record.Status == "running" {
		record.Status = "exited"
		if err := m.save(record); err != nil {
			return false, fmt.Errorf("record exited process: %w", err)
		}
	}
	return running, nil
}

func (m *Manager) Logs(id string, output io.Writer) error {
	record, err := m.Get(id)
	if err != nil {
		return err
	}
	file, err := os.Open(record.LogPath)
	if err != nil {
		return fmt.Errorf("open process log: %w", err)
	}
	_, readErr := io.Copy(output, file)
	closeErr := file.Close()
	if readErr != nil {
		return fmt.Errorf("read process log: %w", readErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close process log: %w", closeErr)
	}
	return nil
}

func processEnvironment(spec Spec) ([]string, error) {
	values := make(map[string]string)
	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			values[key] = value
		}
	}
	for key, value := range spec.Environment {
		values[key] = value
	}
	if len(spec.Secrets) > 0 {
		secrets, err := bitpacksecrets.ResolveFromEnvironment(spec.Secrets)
		if err != nil {
			return nil, err
		}
		for name, value := range secrets {
			values[bitpacksecrets.EnvironmentKey(name)] = value
		}
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	environment := make([]string, 0, len(keys))
	for _, key := range keys {
		environment = append(environment, key+"="+values[key])
	}
	return environment, nil
}

func (m *Manager) save(record *Record) error {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("encode process state: %w", err)
	}
	temporary, err := os.CreateTemp(m.Directory, ".process-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary process state: %w", err)
	}
	name := temporary.Name()
	defer os.Remove(name)
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("write process state: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close process state: %w", err)
	}
	if err := os.Rename(name, m.statePath(record.ID)); err != nil {
		return fmt.Errorf("save process state: %w", err)
	}
	return nil
}

func (m *Manager) statePath(id string) string {
	return filepath.Join(m.Directory, id+".json")
}

func ShellCommand(line string) (string, []string, error) {
	if strings.TrimSpace(line) == "" {
		return "", nil, fmt.Errorf("process command cannot be empty")
	}
	if os.PathSeparator == '\\' {
		executable, err := exec.LookPath("cmd.exe")
		if err != nil {
			return "", nil, fmt.Errorf("find cmd.exe: %w", err)
		}
		return executable, []string{"/D", "/S", "/C", line}, nil
	}
	executable, err := exec.LookPath("sh")
	if err != nil {
		return "", nil, fmt.Errorf("find sh: %w", err)
	}
	return executable, []string{"-c", line}, nil
}

func IDPart(value string) string {
	var result strings.Builder
	for _, char := range strings.ToLower(value) {
		if char >= 'a' && char <= 'z' || char >= '0' && char <= '9' || char == '-' || char == '_' {
			result.WriteRune(char)
		} else if result.Len() > 0 {
			result.WriteByte('-')
		}
	}
	return strings.Trim(result.String(), "-_")
}

func IDWithIndex(prefix string, index int) string {
	return prefix + "-" + strconv.Itoa(index)
}
