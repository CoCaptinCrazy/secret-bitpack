package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Parse(source io.Reader) (*Environment, error) {
	environment := &Environment{
		Workspace: "./",
		Variables: make(map[string]string),
	}
	scanner := bufio.NewScanner(source)
	section := ""
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if section != "" {
			if line == "|:" {
				section = ""
				continue
			}
			switch section {
			case "env":
				key, value, ok := strings.Cut(line, "=")
				if !ok || strings.TrimSpace(key) == "" {
					return nil, fmt.Errorf("line %d: expected KEY = VALUE in env section", lineNumber)
				}
				environment.Variables[strings.TrimSpace(key)] = unquote(strings.TrimSpace(value))
			case "startup":
				environment.Startup = append(environment.Startup, line)
			case "shell":
				environment.Shell = append(environment.Shell, line)
			case "secrets":
				environment.Secrets = append(environment.Secrets, line)
			}
			continue
		}

		switch {
		case strings.HasPrefix(line, "development "):
			parts := strings.SplitN(line, "\"", 3)
			if len(parts) != 3 || strings.TrimSpace(parts[0]) != "development" {
				return nil, fmt.Errorf("line %d: expected development \"name\" v<version>", lineNumber)
			}
			environment.Name = parts[1]
			version := strings.TrimSpace(parts[2])
			if version != "" {
				environment.Version = strings.TrimPrefix(version, "v")
			}
		case strings.HasPrefix(line, "using "):
			environment.Runtime = strings.TrimSpace(strings.TrimPrefix(line, "using "))
		case strings.HasPrefix(line, "workspace "):
			environment.Workspace = unquote(strings.TrimSpace(strings.TrimPrefix(line, "workspace ")))
		case strings.HasPrefix(line, "service "):
			declaration := strings.TrimSpace(strings.TrimPrefix(line, "service "))
			name, command, hasCommand := strings.Cut(declaration, "=")
			name = strings.TrimSpace(name)
			if name == "" {
				return nil, fmt.Errorf("line %d: service name is required", lineNumber)
			}
			service := Service{Name: name}
			if hasCommand {
				service.Command = unquote(strings.TrimSpace(command))
				if service.Command == "" {
					return nil, fmt.Errorf("line %d: service command is required after '='", lineNumber)
				}
			}
			environment.Services = append(environment.Services, service)
		case strings.HasPrefix(line, "port "):
			port, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "port ")))
			if err != nil || port < 1 || port > 65535 {
				return nil, fmt.Errorf("line %d: port must be between 1 and 65535", lineNumber)
			}
			environment.Ports = append(environment.Ports, port)
		case strings.HasPrefix(line, "watch "):
			environment.Watch = append(environment.Watch,
				unquote(strings.TrimSpace(strings.TrimPrefix(line, "watch "))))
		case strings.HasPrefix(line, "env") && strings.HasSuffix(line, ":|"):
			section = "env"
		case strings.HasPrefix(line, "startup") && strings.HasSuffix(line, ":|"):
			section = "startup"
		case strings.HasPrefix(line, "shell") && strings.HasSuffix(line, ":|"):
			section = "shell"
		case strings.HasPrefix(line, "secrets") && strings.HasSuffix(line, ":|"):
			section = "secrets"
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if section != "" {
		return nil, fmt.Errorf("unclosed %s section", section)
	}
	if environment.Name == "" {
		return nil, fmt.Errorf("development declaration not found")
	}
	return environment, nil
}

func unquote(value string) string {
	return strings.Trim(value, "\"'")
}
