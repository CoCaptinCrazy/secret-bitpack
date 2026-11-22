package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Parse(source io.Reader) (*Plan, error) {
	plan := &Plan{Version: "1.0", FailFast: true}
	section := ""
	scanner := bufio.NewScanner(source)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "batch ") {
			parts := strings.SplitN(line, "\"", 3)
			if len(parts) != 3 {
				return nil, fmt.Errorf("line %d: expected batch \"name\" v<version>", lineNumber)
			}
			plan.Name = parts[1]
			version := strings.Fields(strings.TrimSpace(parts[2]))
			if len(version) != 0 {
				plan.Version = strings.TrimPrefix(version[0], "v")
			}
			continue
		}
		key, value, hasValue := strings.Cut(line, ":")
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if hasValue && value != "" {
			switch key {
			case "author":
				plan.Author = value
			case "parallel":
				plan.Parallel = strings.EqualFold(value, "true")
			case "failfast":
				plan.FailFast = strings.EqualFold(value, "true")
			default:
				return nil, fmt.Errorf("line %d: unknown batch setting %q", lineNumber, key)
			}
			continue
		}
		if hasValue {
			switch key {
			case "author", "requires", "tags", "forge", "verify", "publish", "assemble", "ignite":
				section = key
				continue
			case "parallel", "failfast":
				section = key
				continue
			default:
				return nil, fmt.Errorf("line %d: unknown batch section %q", lineNumber, key)
			}
		}
		if section == "" {
			return nil, fmt.Errorf("line %d: expected a batch section", lineNumber)
		}
		switch section {
		case "author":
			plan.Author = line
			section = ""
		case "parallel":
			plan.Parallel = strings.EqualFold(line, "true")
			section = ""
		case "failfast":
			plan.FailFast = strings.EqualFold(line, "true")
			section = ""
		case "requires":
			plan.Requires = append(plan.Requires, line)
		case "tags":
			plan.Tags = append(plan.Tags, line)
		case "forge":
			plan.Forge = append(plan.Forge, line)
		case "verify":
			plan.Verify = append(plan.Verify, line)
		case "publish":
			plan.Publish = append(plan.Publish, line)
		case "assemble":
			plan.Assemble = append(plan.Assemble, line)
		case "ignite":
			plan.Ignite = append(plan.Ignite, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if plan.Name == "" {
		return nil, fmt.Errorf("batch declaration not found")
	}
	return plan, nil
}
