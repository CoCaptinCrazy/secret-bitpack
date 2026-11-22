package kubernetes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(context.Context, string, []string, io.Reader, io.Writer, io.Writer) error
}

type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, executable string, args []string, input io.Reader, output, errorOutput io.Writer) error {
	command := exec.CommandContext(ctx, executable, args...)
	command.Stdin = input
	command.Stdout = output
	command.Stderr = errorOutput
	return command.Run()
}

type Client struct {
	Executable string
	Namespace  string
	Runner     Runner
}

func (c Client) Apply(ctx context.Context, objects []map[string]any, secretValues map[string]string, output io.Writer) error {
	if len(objects) == 0 {
		return fmt.Errorf("no Kubernetes resources to apply")
	}
	var manifest bytes.Buffer
	for index, object := range objects {
		if index > 0 {
			manifest.WriteString("\n---\n")
		}
		encoder := json.NewEncoder(&manifest)
		if err := encoder.Encode(object); err != nil {
			return fmt.Errorf("encode Kubernetes resource: %w", err)
		}
	}
	return c.applyManifest(ctx, manifest.Bytes(), secretValues, output)
}

func (c Client) ApplyAutoscaler(ctx context.Context, assembly, module string, min, max, cpu int, output io.Writer) error {
	assemblyName := resourceName(assembly)
	name := resourceName(module)
	if assemblyName == "" || name == "" || min < 1 || max < min || cpu < 1 || cpu > 100 {
		return fmt.Errorf("invalid autoscaler configuration for module %q", module)
	}
	object := map[string]any{
		"apiVersion": "autoscaling/v2",
		"kind":       "HorizontalPodAutoscaler",
		"metadata": map[string]any{
			"name": name,
			"labels": map[string]string{
				"bitpack.io/assembly": assemblyName,
				"bitpack.io/module":   name,
			},
		},
		"spec": map[string]any{
			"scaleTargetRef": map[string]string{"apiVersion": "apps/v1", "kind": "Deployment", "name": name},
			"minReplicas":    min,
			"maxReplicas":    max,
			"metrics": []any{map[string]any{
				"type": "Resource",
				"resource": map[string]any{
					"name": "cpu",
					"target": map[string]any{
						"type": "Utilization", "averageUtilization": cpu,
					},
				},
			}},
		},
	}
	data, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("encode Kubernetes autoscaler: %w", err)
	}
	return c.applyManifest(ctx, data, nil, output)
}

func (c Client) applyManifest(ctx context.Context, manifest []byte, secretValues map[string]string, output io.Writer) error {
	args := c.namespaceArgs()
	args = append(args, "apply", "-f", "-")
	stdout, stderr, err := c.run(ctx, args, bytes.NewReader(manifest))
	if writeErr := writeOutput(output, redact(stdout, secretValues)); writeErr != nil {
		return writeErr
	}
	if err != nil {
		message := strings.TrimSpace(redact(stderr, secretValues))
		if message == "" {
			return fmt.Errorf("kubectl apply failed using the current kubeconfig context: %w", err)
		}
		return fmt.Errorf("kubectl apply failed using the current kubeconfig context: %s: %w", message, err)
	}
	if warning := redact(stderr, secretValues); warning != "" {
		if err := writeOutput(output, warning); err != nil {
			return err
		}
	}
	return nil
}

func (c Client) Scale(ctx context.Context, deployment string, replicas int, output io.Writer) error {
	if replicas < 0 {
		return fmt.Errorf("replicas must be non-negative")
	}
	name := resourceName(deployment)
	if name == "" {
		return fmt.Errorf("invalid deployment name %q", deployment)
	}
	args := append(c.namespaceArgs(), "scale", "deployment/"+name, fmt.Sprintf("--replicas=%d", replicas))
	return c.runAndReport(ctx, args, output, "scale deployment")
}

func (c Client) Stop(ctx context.Context, deployment string, output io.Writer) error {
	name := resourceName(deployment)
	if name == "" {
		return fmt.Errorf("invalid deployment name %q", deployment)
	}
	deleteArgs := append(c.namespaceArgs(), "delete", "hpa/"+name, "--ignore-not-found=true")
	if err := c.runAndReport(ctx, deleteArgs, output, "disable deployment autoscaler"); err != nil {
		return err
	}
	return c.Scale(ctx, name, 0, output)
}

func (c Client) Restart(ctx context.Context, deployment string, output io.Writer) error {
	name := resourceName(deployment)
	if name == "" {
		return fmt.Errorf("invalid deployment name %q", deployment)
	}
	args := append(c.namespaceArgs(), "rollout", "restart", "deployment/"+name)
	return c.runAndReport(ctx, args, output, "restart deployment")
}

func (c Client) Rollback(ctx context.Context, deployment string, output io.Writer) error {
	name := resourceName(deployment)
	if name == "" {
		return fmt.Errorf("invalid deployment name %q", deployment)
	}
	args := append(c.namespaceArgs(), "rollout", "undo", "deployment/"+name)
	return c.runAndReport(ctx, args, output, "rollback deployment")
}

func (c Client) DeleteModule(ctx context.Context, module string, output io.Writer) error {
	name := resourceName(module)
	if name == "" {
		return fmt.Errorf("invalid module name %q", module)
	}
	args := append(c.namespaceArgs(), "delete", "deployment,service,pvc,networkpolicy",
		"--selector", "bitpack.io/module="+name, "--ignore-not-found=true")
	return c.runAndReport(ctx, args, output, "delete module resources")
}

func (c Client) DeleteAssembly(ctx context.Context, assembly string, output io.Writer) error {
	name := resourceName(assembly)
	if name == "" {
		return fmt.Errorf("invalid assembly name %q", assembly)
	}
	selector := "bitpack.io/assembly=" + name
	args := append(c.namespaceArgs(), "delete", "all,hpa,networkpolicy,pvc,secrets",
		"--selector", selector, "--ignore-not-found=true")
	return c.runAndReport(ctx, args, output, "delete assembly resources")
}

func (c Client) Logs(ctx context.Context, deployment string, output io.Writer) error {
	name := resourceName(deployment)
	if name == "" {
		return fmt.Errorf("invalid deployment name %q", deployment)
	}
	args := append(c.namespaceArgs(), "logs", "deployment/"+name, "--all-containers=true")
	return c.runAndReport(ctx, args, output, "read deployment logs")
}

func (c Client) Health(ctx context.Context, deployment string, output io.Writer) error {
	name := resourceName(deployment)
	if name == "" {
		return fmt.Errorf("invalid deployment name %q", deployment)
	}
	args := append(c.namespaceArgs(), "rollout", "status", "deployment/"+name, "--timeout=10s")
	return c.runAndReport(ctx, args, output, "check deployment health")
}

func (c Client) namespaceArgs() []string {
	if c.Namespace == "" {
		return nil
	}
	return []string{"--namespace", c.Namespace}
}

func (c Client) runAndReport(ctx context.Context, args []string, output io.Writer, action string) error {
	stdout, stderr, err := c.run(ctx, args, nil)
	if writeErr := writeOutput(output, stdout); writeErr != nil {
		return writeErr
	}
	if err != nil {
		message := strings.TrimSpace(stderr)
		if message == "" {
			return fmt.Errorf("kubectl %s failed using the current kubeconfig context: %w", action, err)
		}
		return fmt.Errorf("kubectl %s failed using the current kubeconfig context: %s: %w", action, message, err)
	}
	return writeOutput(output, stderr)
}

func (c Client) run(ctx context.Context, args []string, input io.Reader) (string, string, error) {
	executable := c.Executable
	if executable == "" {
		executable = "kubectl"
	}
	runner := c.Runner
	if runner == nil {
		runner = ExecRunner{}
	}
	var stdout, stderr bytes.Buffer
	err := runner.Run(ctx, executable, args, input, &stdout, &stderr)
	return stdout.String(), stderr.String(), err
}

func redact(value string, secrets map[string]string) string {
	for _, secret := range secrets {
		if secret != "" {
			value = strings.ReplaceAll(value, secret, "[REDACTED]")
		}
	}
	return value
}

func writeOutput(output io.Writer, value string) error {
	if value == "" {
		return nil
	}
	if _, err := io.WriteString(output, value); err != nil {
		return fmt.Errorf("write kubectl output: %w", err)
	}
	return nil
}
