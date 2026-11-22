package kubernetes

import (
	"fmt"
	"sort"
	"strings"

	assemblyparser "secret-bitpack/internal/assembly/parser"
	bitpacksecrets "secret-bitpack/internal/security/secrets"
)

func Compile(assembly *assemblyparser.BPack, secretValues map[string]string) ([]map[string]any, error) {
	if assembly == nil || assembly.Name == "" {
		return nil, fmt.Errorf("assembly is required")
	}
	assemblyName := resourceName(assembly.Name)
	if assemblyName == "" {
		return nil, fmt.Errorf("assembly name %q is not a valid Kubernetes resource name", assembly.Name)
	}
	scaled := make(map[string]assemblyparser.Autoscale, len(assembly.Autoscale))
	for _, scale := range assembly.Autoscale {
		if _, exists := scaled[scale.Target]; exists {
			return nil, fmt.Errorf("module %q has multiple autoscale policies", scale.Target)
		}
		scaled[scale.Target] = scale
	}
	policyAnnotations := make(map[string]string)
	for _, policy := range assembly.Policies {
		if policy.Restart != "" && policy.Restart != "always" {
			return nil, fmt.Errorf("Kubernetes Deployments require restart always; policy %q requests %q",
				policy.Name, policy.Restart)
		}
		if policy.Recovery != "" && policy.Recovery != "automatic" {
			return nil, fmt.Errorf("Kubernetes Deployments use automatic recovery; policy %q requests %q",
				policy.Name, policy.Recovery)
		}
		policyAnnotations["bitpack.io/policy-"+resourceName(policy.Name)] =
			"restart=always,recovery=automatic"
	}

	objects := make([]map[string]any, 0, len(assembly.Modules)*3+len(assembly.Autoscale)+len(assembly.Rules))
	secretValuesUsed := make(map[string]string)
	secretName := resourceName(assemblyName + "-secrets")
	modules := append([]assemblyparser.Module(nil), assembly.Modules...)
	sort.Slice(modules, func(i, j int) bool { return modules[i].Name < modules[j].Name })

	for _, module := range modules {
		name := resourceName(module.Name)
		if name == "" {
			return nil, fmt.Errorf("module name %q is not a valid Kubernetes resource name", module.Name)
		}
		labels := map[string]string{
			"app.kubernetes.io/name":       name,
			"app.kubernetes.io/managed-by": "bitpack",
			"bitpack.io/assembly":          assemblyName,
			"bitpack.io/module":            name,
		}
		if module.Mesh != "" {
			labels["bitpack.io/mesh"] = resourceName(module.Mesh)
		}
		annotations := make(map[string]string, len(policyAnnotations)+1)
		for key, value := range policyAnnotations {
			annotations[key] = value
		}
		if len(module.Depends) > 0 {
			dependencies := append([]string(nil), module.Depends...)
			sort.Strings(dependencies)
			annotations["bitpack.io/depends-on"] = strings.Join(dependencies, ",")
		}
		selector := map[string]any{"matchLabels": map[string]string{"app.kubernetes.io/name": name}}

		container := map[string]any{
			"name":  name,
			"image": module.Blueprint,
		}
		containerPorts := make([]any, 0, len(module.Bridges))
		servicePorts := make([]any, 0, len(module.Bridges))
		seenPortNames := make(map[string]bool)
		for _, bridge := range module.Bridges {
			portName := resourceName(bridge.Name)
			if portName == "" || seenPortNames[portName] {
				return nil, fmt.Errorf("module %q has an invalid or duplicate bridge name %q", module.Name, bridge.Name)
			}
			seenPortNames[portName] = true
			containerPorts = append(containerPorts, map[string]any{
				"name": portName, "containerPort": bridge.TargetPort, "protocol": "TCP",
			})
			servicePorts = append(servicePorts, map[string]any{
				"name": portName, "port": bridge.HostPort, "targetPort": bridge.TargetPort, "protocol": "TCP",
			})
		}
		if len(containerPorts) > 0 {
			container["ports"] = containerPorts
		}

		environment := make([]any, 0, len(module.Env)+len(module.Secrets))
		envNames := make(map[string]bool, len(module.Env)+len(module.Secrets))
		envKeys := make([]string, 0, len(module.Env))
		for key := range module.Env {
			envKeys = append(envKeys, key)
		}
		sort.Strings(envKeys)
		for _, key := range envKeys {
			if !validEnvironmentName(key) {
				return nil, fmt.Errorf("module %q has invalid environment variable %q", module.Name, key)
			}
			envNames[key] = true
			environment = append(environment, map[string]any{"name": key, "value": module.Env[key]})
		}
		for _, secret := range module.Secrets {
			value, exists := secretValues[secret]
			if !exists {
				return nil, fmt.Errorf("module %q requires unresolved secret %q", module.Name, secret)
			}
			envName := bitpacksecrets.EnvironmentKey(secret)
			if envNames[envName] {
				return nil, fmt.Errorf("module %q has duplicate environment variable %q", module.Name, envName)
			}
			envNames[envName] = true
			secretValuesUsed[secret] = value
			environment = append(environment, map[string]any{
				"name": envName,
				"valueFrom": map[string]any{
					"secretKeyRef": map[string]any{"name": secretName, "key": secret},
				},
			})
		}
		if len(environment) > 0 {
			container["env"] = environment
		}

		volumeMounts := make([]any, 0, len(module.Volumes))
		volumes := make([]any, 0, len(module.Volumes))
		for _, volume := range module.Volumes {
			volumeName := resourceName(name + "-" + volume.Name)
			if volumeName == "" {
				return nil, fmt.Errorf("module %q has invalid volume name %q", module.Name, volume.Name)
			}
			volumeMounts = append(volumeMounts, map[string]any{
				"name": volumeName, "mountPath": volume.Path,
			})
			volumes = append(volumes, map[string]any{
				"name": volumeName, "persistentVolumeClaim": map[string]any{"claimName": volumeName},
			})
			objects = append(objects, object("v1", "PersistentVolumeClaim", volumeName, labels,
				map[string]any{"accessModes": []string{"ReadWriteOnce"}, "resources": map[string]any{
					"requests": map[string]string{"storage": volume.Size},
				}},
			))
		}
		if len(volumeMounts) > 0 {
			container["volumeMounts"] = volumeMounts
		}
		podSpec := map[string]any{
			"containers":    []any{container},
			"restartPolicy": "Always",
		}
		if len(volumes) > 0 {
			podSpec["volumes"] = volumes
		}
		if _, hasAutoscale := scaled[module.Name]; hasAutoscale {
			container["resources"] = map[string]any{
				"requests": map[string]string{"cpu": "100m"},
			}
		}
		podTemplate := map[string]any{
			"metadata": map[string]any{"labels": map[string]string{
				"app.kubernetes.io/name": name, "bitpack.io/assembly": assemblyName,
			}},
			"spec": podSpec,
		}
		deployment := object("apps/v1", "Deployment", name, labels, map[string]any{
			"replicas": module.Replicas,
			"selector": selector,
			"template": podTemplate,
		})
		if len(annotations) > 0 {
			deployment["metadata"].(map[string]any)["annotations"] = annotations
		}
		objects = append(objects, deployment)
		objects = append(objects, object("v1", "Service", name, labels, map[string]any{
			"selector": map[string]string{"app.kubernetes.io/name": name},
			"ports":    servicePorts,
		}))
	}

	for _, scale := range assembly.Autoscale {
		name := resourceName(scale.Target)
		objects = append(objects, object("autoscaling/v2", "HorizontalPodAutoscaler",
			name, map[string]string{
				"bitpack.io/assembly": assemblyName,
				"bitpack.io/module":   name,
			}, map[string]any{
				"scaleTargetRef": map[string]string{"apiVersion": "apps/v1", "kind": "Deployment", "name": name},
				"minReplicas":    scale.Min,
				"maxReplicas":    scale.Max,
				"metrics": []any{map[string]any{
					"type": "Resource",
					"resource": map[string]any{
						"name": "cpu",
						"target": map[string]any{
							"type": "Utilization", "averageUtilization": scale.CPU,
						},
					},
				}},
			}))
	}

	if len(assembly.Rules) > 0 {
		rulesByDestination := make(map[string][]map[string]any)
		destinations := make(map[string]bool)
		for _, rule := range assembly.Rules {
			destinations[rule.Destination] = true
			if rule.Action == "allow" {
				rulesByDestination[rule.Destination] = append(
					rulesByDestination[rule.Destination],
					map[string]any{"from": []any{map[string]any{
						"podSelector": map[string]any{"matchLabels": map[string]string{
							"app.kubernetes.io/name": resourceName(rule.Source),
						}},
					}}},
				)
			}
		}
		for destination := range destinations {
			name := resourceName(destination)
			objects = append(objects, object("networking.k8s.io/v1", "NetworkPolicy", name,
				map[string]string{"bitpack.io/assembly": assemblyName}, map[string]any{
					"podSelector": map[string]any{"matchLabels": map[string]string{
						"app.kubernetes.io/name": name,
					}},
					"policyTypes": []string{"Ingress"},
					"ingress":     rulesByDestination[destination],
				}))
		}
	}
	if len(secretValuesUsed) > 0 {
		objects = append(objects, map[string]any{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]any{
				"name": secretName,
				"labels": map[string]string{
					"app.kubernetes.io/managed-by": "bitpack",
					"bitpack.io/assembly":          assemblyName,
				},
			},
			"type":       "Opaque",
			"stringData": secretValuesUsed,
		})
	}
	return objects, nil
}

func object(apiVersion, kind, name string, labels map[string]string, spec map[string]any) map[string]any {
	return map[string]any{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]any{
			"name": name, "labels": labels,
		},
		"spec": spec,
	}
}

func resourceName(value string) string {
	var result strings.Builder
	lastDash := false
	for _, char := range strings.ToLower(value) {
		if char >= 'a' && char <= 'z' || char >= '0' && char <= '9' {
			result.WriteRune(char)
			lastDash = false
		} else if !lastDash && result.Len() > 0 {
			result.WriteByte('-')
			lastDash = true
		}
		if result.Len() >= 63 {
			break
		}
	}
	return strings.Trim(result.String(), "-")
}

func validEnvironmentName(value string) bool {
	if value == "" {
		return false
	}
	for index, char := range value {
		if char == '_' || char >= 'A' && char <= 'Z' ||
			index > 0 && char >= '0' && char <= '9' {
			continue
		}
		return false
	}
	return true
}
