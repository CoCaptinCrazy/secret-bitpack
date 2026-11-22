package kubernetes

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	assemblyparser "github.com/CoCaptinCrazy/secret-bitpack/internal/assembly/parser"
)

func TestCompileExampleBPack(t *testing.T) {
	path := filepath.Join("..", "..", "..", "examples", "bitpack.bpack")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	assembly, err := assemblyparser.ParseBPack(file)
	if err != nil {
		t.Fatal(err)
	}
	objects, err := Compile(assembly, nil)
	if err != nil {
		t.Fatal(err)
	}
	counts := make(map[string]int)
	for _, resource := range objects {
		counts[resource["kind"].(string)]++
	}
	if counts["Deployment"] != 6 || counts["Service"] != 6 ||
		counts["PersistentVolumeClaim"] != 3 || counts["HorizontalPodAutoscaler"] != 2 ||
		counts["NetworkPolicy"] != 3 {
		t.Fatalf("unexpected generated resources: %#v", counts)
	}
	api := findObject(objects, "Deployment", "api")
	if api == nil {
		t.Fatal("API Deployment was not generated")
	}
	spec := api["spec"].(map[string]any)
	template := spec["template"].(map[string]any)
	podSpec := template["spec"].(map[string]any)
	containers := podSpec["containers"].([]any)
	container := containers[0].(map[string]any)
	if container["image"] != "social-api" {
		t.Fatalf("blueprint image reference was not retained: %#v", container["image"])
	}
}

func TestCompileInjectsOnlySecretReferencesIntoPods(t *testing.T) {
	assembly := &assemblyparser.BPack{
		Name: "demo",
		Modules: []assemblyparser.Module{{
			Name: "api", Blueprint: "example/api:1", Replicas: 1,
			Secrets: []string{"db-password"},
		}},
	}
	values := map[string]string{"db-password": "test-secret-value"}
	objects, err := Compile(assembly, values)
	if err != nil {
		t.Fatal(err)
	}
	deployment := findObject(objects, "Deployment", "api")
	container := deployment["spec"].(map[string]any)["template"].(map[string]any)["spec"].(map[string]any)["containers"].([]any)[0].(map[string]any)
	env := container["env"].([]any)[0].(map[string]any)
	if env["name"] != "BITPACK_SECRET_DB_PASSWORD" {
		t.Fatalf("unexpected secret environment name: %#v", env["name"])
	}
	valueFrom := env["valueFrom"].(map[string]any)
	if valueFrom["secretKeyRef"] == nil {
		t.Fatalf("pod contains no Kubernetes Secret reference: %#v", env)
	}
	deploymentJSON, err := json.Marshal(deployment)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(deploymentJSON), "test-secret-value") {
		t.Fatal("secret value was copied into a pod environment literal")
	}
	secret := findObject(objects, "Secret", "demo-secrets")
	if secret == nil || secret["stringData"].(map[string]string)["db-password"] != "test-secret-value" {
		t.Fatalf("Kubernetes Secret object was not generated: %#v", secret)
	}
}

func TestCompileRejectsUnresolvedSecret(t *testing.T) {
	assembly := &assemblyparser.BPack{
		Name: "demo",
		Modules: []assemblyparser.Module{{
			Name: "api", Blueprint: "example/api:1", Replicas: 1,
			Secrets: []string{"db-password"},
		}},
	}
	if _, err := Compile(assembly, nil); err == nil || !strings.Contains(err.Error(), "unresolved secret") {
		t.Fatalf("expected unresolved secret error, got %v", err)
	}
}

func findObject(objects []map[string]any, kind, name string) map[string]any {
	for _, object := range objects {
		if object["kind"] != kind {
			continue
		}
		metadata := object["metadata"].(map[string]any)
		if metadata["name"] == name {
			return object
		}
	}
	return nil
}
