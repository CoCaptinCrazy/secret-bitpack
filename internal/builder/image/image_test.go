package image

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type fakeRegistry struct {
	base        v1.Image
	pulled      string
	pushedTag   string
	pushedImage v1.Image
}

func (r *fakeRegistry) Pull(_ context.Context, reference string) (v1.Image, error) {
	r.pulled = reference
	return r.base, nil
}

func (r *fakeRegistry) Push(_ context.Context, reference string, image v1.Image) error {
	r.pushedTag = reference
	r.pushedImage = image
	return nil
}

func TestBuildStoresAndPushesNativeOCILayout(t *testing.T) {
	workspace := t.TempDir()
	for name, content := range map[string]string{
		"Bitfile": "blueprint \"demo\" v1\n",
		"app.txt": "application payload",
		".env":    "SECRET=not-in-image",
	} {
		if err := os.WriteFile(filepath.Join(workspace, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	registry := &fakeRegistry{base: empty.Image}
	builder := Builder{
		Registry: registry, StoreDirectory: filepath.Join(t.TempDir(), "images"),
	}
	tag := "registry.example/team/demo:v1"
	digest, err := builder.Build(context.Background(), Recipe{
		BaseImage: "node:22-alpine", Tag: tag, Name: "demo", Workspace: workspace,
		Environment: map[string]string{"APP_ENV": "test"},
		Command:     []string{"node server.js"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if digest == "" || registry.pulled != "node:22-alpine" {
		t.Fatalf("base image or result digest missing: digest=%q pulled=%q", digest, registry.pulled)
	}
	if err := builder.Push(context.Background(), tag); err != nil {
		t.Fatal(err)
	}
	if registry.pushedTag != tag || registry.pushedImage == nil {
		t.Fatal("local OCI image was not passed to the registry")
	}
	mediaType, err := registry.pushedImage.MediaType()
	if err != nil {
		t.Fatal(err)
	}
	if mediaType != types.OCIManifestSchema1 {
		t.Fatalf("expected OCI manifest media type, got %s", mediaType)
	}
	config, err := registry.pushedImage.ConfigFile()
	if err != nil {
		t.Fatal(err)
	}
	if config.Config.WorkingDir != "/app" || config.Config.Entrypoint[0] != "/bin/sh" {
		t.Fatalf("unexpected image runtime config: %#v", config.Config)
	}
	foundEnvironment := false
	for _, value := range config.Config.Env {
		if value == "APP_ENV=test" {
			foundEnvironment = true
		}
	}
	if !foundEnvironment {
		t.Fatalf("image environment was not applied: %#v", config.Config.Env)
	}

	layers, err := registry.pushedImage.Layers()
	if err != nil {
		t.Fatal(err)
	}
	if len(layers) != 1 {
		t.Fatalf("expected one project layer, got %d", len(layers))
	}
	layerReader, err := layers[0].Uncompressed()
	if err != nil {
		layoutRoot := filepath.Join(builder.StoreDirectory, referenceKey(tag))
		var files []string
		_ = filepath.WalkDir(layoutRoot, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr == nil && !entry.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		t.Fatalf("%v; OCI files: %v", err, files)
	}
	defer layerReader.Close()
	archive := tar.NewReader(layerReader)
	foundPayload := false
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if header.Name == "app/app.txt" {
			payload, err := io.ReadAll(archive)
			if err != nil {
				t.Fatal(err)
			}
			foundPayload = string(payload) == "application payload"
		}
		if header.Name == "app/.env" {
			t.Fatal("environment secret file was included in the image")
		}
	}
	if !foundPayload {
		t.Fatal("project files were not included in the OCI layer")
	}
}
