package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/builder/blueprints"
)

func Parse(r io.Reader) (*blueprints.Blueprint, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	var blueprint *blueprints.Blueprint
	section := ""
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "\ufeff"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if section != "" {
			if line == "|>" {
				section = ""
				continue
			}
			if err := addSectionValue(blueprint, section, line, lineNumber); err != nil {
				return nil, err
			}
			continue
		}
		if strings.HasPrefix(line, "blueprint ") {
			if blueprint != nil {
				return nil, fmt.Errorf("line %d: duplicate blueprint declaration", lineNumber)
			}
			name, version, err := parseDeclaration(line)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNumber, err)
			}
			blueprint = blueprints.New(name, version)
			continue
		}
		if blueprint == nil {
			return nil, fmt.Errorf("line %d: blueprint declaration must come first", lineNumber)
		}
		if strings.HasSuffix(line, "<|") {
			sectionName := strings.TrimSpace(strings.TrimSuffix(line, "<|"))
			if sectionName == "" {
				return nil, fmt.Errorf("line %d: section name is required", lineNumber)
			}
			section = sectionName
			continue
		}
		boundary := strings.IndexAny(line, " \t")
		if boundary < 1 {
			return nil, fmt.Errorf("line %d: expected a directive and value", lineNumber)
		}
		key := line[:boundary]
		value := unquote(strings.TrimSpace(line[boundary:]))
		if value == "" {
			return nil, fmt.Errorf("line %d: expected a directive and value", lineNumber)
		}
		blueprint.Directives[key] = value
		switch key {
		case "using":
			blueprint.Using = value
		case "bring":
			blueprint.Context = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if section != "" {
		return nil, fmt.Errorf("unclosed %q section", section)
	}
	if blueprint == nil {
		return nil, fmt.Errorf("blueprint declaration not found")
	}
	return blueprint, nil
}

func parseDeclaration(line string) (string, string, error) {
	parts := strings.SplitN(line, "\"", 3)
	if len(parts) != 3 || strings.TrimSpace(parts[0]) != "blueprint" {
		return "", "", fmt.Errorf("expected blueprint \"name\" [version]")
	}
	name := strings.TrimSpace(parts[1])
	if name == "" {
		return "", "", fmt.Errorf("blueprint name is required")
	}
	version := strings.TrimSpace(parts[2])
	if version == "" {
		return name, "1.0", nil
	}
	fields := strings.Fields(version)
	if len(fields) != 1 {
		return "", "", fmt.Errorf("invalid blueprint version")
	}
	return name, strings.TrimPrefix(fields[0], "v"), nil
}

func addSectionValue(blueprint *blueprints.Blueprint, section, line string, lineNumber int) error {
	values := blueprint.Sections[section]
	switch section {
	case "env":
		key, value, ok := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return fmt.Errorf("line %d: expected KEY = VALUE in env section", lineNumber)
		}
		value = unquote(strings.TrimSpace(value))
		blueprint.Environment[key] = value
		blueprint.Sections[section] = append(values, key+"="+value)
	case "secret":
		blueprint.Secrets = append(blueprint.Secrets, line)
		blueprint.Sections[section] = append(values, line)
	case "health":
		blueprint.Health = append(blueprint.Health, line)
		blueprint.Sections[section] = append(values, line)
	case "assemble":
		blueprint.Assemble = append(blueprint.Assemble, line)
		blueprint.Sections[section] = append(values, line)
	case "verify":
		blueprint.Verify = append(blueprint.Verify, line)
		blueprint.Sections[section] = append(values, line)
	case "ignite":
		blueprint.Ignite = append(blueprint.Ignite, line)
		blueprint.Sections[section] = append(values, line)
	default:
		blueprint.Sections[section] = append(values, line)
	}
	return nil
}

func unquote(value string) string {
	return strings.Trim(value, "\"'")
}
