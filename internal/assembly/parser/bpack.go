package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type BPack struct {
	Name      string
	Meshes    []Mesh
	Modules   []Module
	Rules     []FirewallRule
	Autoscale []Autoscale
	Policies  []Policy
}

type Mesh struct {
	Name string
}

type Module struct {
	Name      string
	Blueprint string
	Mesh      string
	Replicas  int
	Depends   []string
	Bridges   []Bridge
	Volumes   []Volume
	Secrets   []string
	Env       map[string]string
}

type Bridge struct {
	Name       string
	HostPort   int
	TargetPort int
}

type Volume struct {
	Name string
	Path string
	Size string
}

type FirewallRule struct {
	Source      string
	Destination string
	Action      string
}

type Autoscale struct {
	Target string
	Min    int
	Max    int
	CPU    int
}

type Policy struct {
	Name     string
	Restart  string
	Recovery string
}

type node struct {
	kind       string
	name       string
	line       int
	properties map[string][]string
	children   []*node
}

func ParseBPack(source io.Reader) (*BPack, error) {
	root := &node{kind: "root"}
	stack := []*node{root}
	scanner := bufio.NewScanner(source)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(stripComment(scanner.Text()))
		if line == "" {
			continue
		}
		if line == "}" {
			if len(stack) == 1 {
				return nil, fmt.Errorf("line %d: unexpected closing brace", lineNumber)
			}
			stack = stack[:len(stack)-1]
			continue
		}
		if strings.HasSuffix(line, "{") {
			header := strings.TrimSpace(strings.TrimSuffix(line, "{"))
			parts, err := splitWords(header)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNumber, err)
			}
			if len(parts) == 0 {
				return nil, fmt.Errorf("line %d: empty block declaration", lineNumber)
			}
			block := &node{
				kind: parts[0], line: lineNumber, properties: make(map[string][]string),
			}
			if len(parts) > 1 {
				block.name = parts[1]
			}
			if len(parts) > 2 {
				return nil, fmt.Errorf("line %d: too many values in %s declaration", lineNumber, block.kind)
			}
			if !knownBlock(block.kind) {
				return nil, fmt.Errorf("line %d: unknown block %q", lineNumber, block.kind)
			}
			stack[len(stack)-1].children = append(stack[len(stack)-1].children, block)
			stack = append(stack, block)
			continue
		}
		if strings.ContainsAny(line, "{}") {
			return nil, fmt.Errorf("line %d: place block braces on their own line", lineNumber)
		}
		if len(stack) == 1 {
			return nil, fmt.Errorf("line %d: property outside a block", lineNumber)
		}
		parts, err := splitWords(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: expected a property and value", lineNumber)
		}
		key := parts[0]
		value := strings.Join(parts[1:], " ")
		current := stack[len(stack)-1]
		current.properties[key] = append(current.properties[key], value)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read BPack: %w", err)
	}
	if len(stack) != 1 {
		unclosed := stack[len(stack)-1]
		return nil, fmt.Errorf("line %d: unclosed %s block", unclosed.line, unclosed.kind)
	}
	if len(root.children) != 1 || root.children[0].kind != "assembly" {
		return nil, fmt.Errorf("BPack must contain exactly one assembly block")
	}
	return convertAssembly(root.children[0])
}

func convertAssembly(root *node) (*BPack, error) {
	if root.name == "" {
		return nil, fmt.Errorf("line %d: assembly name is required", root.line)
	}
	assembly := &BPack{Name: root.name}
	moduleNames := make(map[string]bool)

	for _, child := range root.children {
		switch child.kind {
		case "mesh":
			if child.name == "" {
				return nil, fmt.Errorf("line %d: mesh name is required", child.line)
			}
			assembly.Meshes = append(assembly.Meshes, Mesh{Name: child.name})
			for _, nested := range child.children {
				if err := addModule(assembly, nested, child.name, moduleNames); err != nil {
					return nil, err
				}
			}
		case "module":
			if err := addModule(assembly, child, "", moduleNames); err != nil {
				return nil, err
			}
		case "firewall":
			for _, ruleNode := range child.children {
				if ruleNode.kind != "rule" {
					return nil, fmt.Errorf("line %d: firewall only accepts rule blocks", ruleNode.line)
				}
				rule := FirewallRule{
					Source:      property(ruleNode, "source"),
					Destination: property(ruleNode, "destination"),
					Action:      strings.ToLower(property(ruleNode, "action")),
				}
				if rule.Source == "" || rule.Destination == "" ||
					(rule.Action != "allow" && rule.Action != "deny") {
					return nil, fmt.Errorf("line %d: firewall rule requires source, destination, and allow/deny action", ruleNode.line)
				}
				assembly.Rules = append(assembly.Rules, rule)
			}
		case "autoscale":
			scale, err := convertAutoscale(child)
			if err != nil {
				return nil, err
			}
			assembly.Autoscale = append(assembly.Autoscale, scale)
		case "policy":
			if child.name == "" {
				return nil, fmt.Errorf("line %d: policy name is required", child.line)
			}
			assembly.Policies = append(assembly.Policies, Policy{
				Name: child.name, Restart: property(child, "restart"),
				Recovery: property(child, "recovery"),
			})
		default:
			return nil, fmt.Errorf("line %d: %s is not valid at assembly scope", child.line, child.kind)
		}
	}
	for _, module := range assembly.Modules {
		for _, dependency := range module.Depends {
			if !moduleNames[dependency] {
				return nil, fmt.Errorf("module %q depends on unknown module %q", module.Name, dependency)
			}
		}
	}
	for _, rule := range assembly.Rules {
		if !moduleNames[rule.Source] || !moduleNames[rule.Destination] {
			return nil, fmt.Errorf("firewall rule references unknown module %q or %q",
				rule.Source, rule.Destination)
		}
	}
	for _, scale := range assembly.Autoscale {
		if !moduleNames[scale.Target] {
			return nil, fmt.Errorf("autoscale target %q is not a module", scale.Target)
		}
	}
	if len(assembly.Modules) == 0 {
		return nil, fmt.Errorf("assembly %q contains no modules", assembly.Name)
	}
	return assembly, nil
}

func addModule(assembly *BPack, block *node, mesh string, names map[string]bool) error {
	if block.kind != "module" {
		return fmt.Errorf("line %d: mesh only accepts module blocks", block.line)
	}
	if block.name == "" {
		return fmt.Errorf("line %d: module name is required", block.line)
	}
	if names[block.name] {
		return fmt.Errorf("line %d: duplicate module %q", block.line, block.name)
	}
	names[block.name] = true
	image := property(block, "image")
	if image == "" {
		image = property(block, "blueprint")
	}
	if image == "" {
		return fmt.Errorf("line %d: module %q requires a blueprint image reference", block.line, block.name)
	}
	replicas := 1
	if value := property(block, "replicas"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			return fmt.Errorf("line %d: module %q replicas must be a non-negative integer", block.line, block.name)
		}
		replicas = parsed
	}
	module := Module{
		Name: block.name, Blueprint: image, Mesh: mesh, Replicas: replicas,
		Depends: append([]string(nil), block.properties["depends"]...),
		Secrets: append([]string(nil), block.properties["secret"]...),
		Env:     make(map[string]string),
	}
	for _, raw := range block.properties["env"] {
		key, value, ok := strings.Cut(raw, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return fmt.Errorf("line %d: module env must use KEY=VALUE syntax", block.line)
		}
		module.Env[key] = strings.TrimSpace(value)
	}
	for _, child := range block.children {
		switch child.kind {
		case "bridge":
			bridge, err := convertBridge(child)
			if err != nil {
				return err
			}
			module.Bridges = append(module.Bridges, bridge)
		case "volume":
			volume, err := convertVolume(child)
			if err != nil {
				return err
			}
			module.Volumes = append(module.Volumes, volume)
		default:
			return fmt.Errorf("line %d: module cannot contain %s blocks", child.line, child.kind)
		}
	}
	assembly.Modules = append(assembly.Modules, module)
	return nil
}

func convertBridge(block *node) (Bridge, error) {
	bridge := Bridge{Name: block.name}
	var err error
	bridge.HostPort, err = integerProperty(block, "host")
	if err != nil {
		return Bridge{}, err
	}
	bridge.TargetPort, err = integerProperty(block, "target")
	if err != nil {
		return Bridge{}, err
	}
	if bridge.HostPort < 1 || bridge.HostPort > 65535 ||
		bridge.TargetPort < 1 || bridge.TargetPort > 65535 {
		return Bridge{}, fmt.Errorf("line %d: bridge ports must be between 1 and 65535", block.line)
	}
	if bridge.Name == "" {
		bridge.Name = fmt.Sprintf("port-%d", bridge.HostPort)
	}
	return bridge, nil
}

func convertVolume(block *node) (Volume, error) {
	volume := Volume{Name: block.name, Path: property(block, "path"), Size: property(block, "size")}
	if volume.Name == "" || volume.Path == "" {
		return Volume{}, fmt.Errorf("line %d: volume requires a name and path", block.line)
	}
	if volume.Size == "" {
		volume.Size = "1Gi"
	}
	return volume, nil
}

func convertAutoscale(block *node) (Autoscale, error) {
	scale := Autoscale{Target: property(block, "target")}
	var err error
	scale.Min, err = integerProperty(block, "min")
	if err != nil {
		return Autoscale{}, err
	}
	scale.Max, err = integerProperty(block, "max")
	if err != nil {
		return Autoscale{}, err
	}
	scale.CPU, err = integerProperty(block, "cpu")
	if err != nil {
		return Autoscale{}, err
	}
	if scale.Target == "" || scale.Min < 1 || scale.Max < scale.Min || scale.CPU < 1 || scale.CPU > 100 {
		return Autoscale{}, fmt.Errorf("line %d: autoscale requires target, min >= 1, max >= min, and cpu 1-100",
			block.line)
	}
	return scale, nil
}

func integerProperty(block *node, key string) (int, error) {
	value := property(block, key)
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("line %d: %s must be an integer", block.line, key)
	}
	return parsed, nil
}

func property(block *node, key string) string {
	values := block.properties[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func knownBlock(kind string) bool {
	switch kind {
	case "assembly", "mesh", "module", "bridge", "volume", "firewall", "rule", "autoscale", "policy":
		return true
	default:
		return false
	}
}

func splitWords(line string) ([]string, error) {
	var words []string
	var current strings.Builder
	var quote rune
	escaped := false
	for _, char := range line {
		if escaped {
			current.WriteRune(char)
			escaped = false
			continue
		}
		if char == '\\' && quote != 0 {
			escaped = true
			continue
		}
		if quote != 0 {
			if char == quote {
				quote = 0
			} else {
				current.WriteRune(char)
			}
			continue
		}
		if char == '"' || char == '\'' {
			quote = char
			continue
		}
		if char == ' ' || char == '\t' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			continue
		}
		current.WriteRune(char)
	}
	if escaped {
		current.WriteRune('\\')
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quoted value")
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}
	return words, nil
}

func stripComment(line string) string {
	var quote rune
	escaped := false
	for index, char := range line {
		if escaped {
			escaped = false
			continue
		}
		if char == '\\' && quote != 0 {
			escaped = true
			continue
		}
		if quote != 0 {
			if char == quote {
				quote = 0
			}
			continue
		}
		if char == '"' || char == '\'' {
			quote = char
		} else if char == '#' {
			return line[:index]
		}
	}
	return line
}
