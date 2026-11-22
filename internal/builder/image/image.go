package image

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/CoCaptinCrazy/secret-bitpack/internal/builder/packaging"
)

const imageRefAnnotation = "org.opencontainers.image.ref.name"

type Recipe struct {
	BaseImage   string
	Tag         string
	Name        string
	Workspace   string
	WorkingDir  string
	Environment map[string]string
	Entrypoint  []string
	Command     []string
}

type Registry interface {
	Pull(context.Context, string) (v1.Image, error)
	Push(context.Context, string, v1.Image) error
}

type Builder struct {
	Registry       Registry
	StoreDirectory string
}

type referenceRecord struct {
	Tag    string `json:"tag"`
	Layout string `json:"layout"`
}

type RemoteRegistry struct {
	username string
	password string
}

func NewRemoteRegistry() *RemoteRegistry {
	return &RemoteRegistry{
		username: os.Getenv("BITPACK_REGISTRY_USERNAME"),
		password: os.Getenv("BITPACK_REGISTRY_PASSWORD"),
	}
}

func (r *RemoteRegistry) Pull(ctx context.Context, reference string) (v1.Image, error) {
	imageRef, err := name.ParseReference(reference, name.WeakValidation)
	if err != nil {
		return nil, fmt.Errorf("parse base image reference %q: %w", reference, err)
	}
	auth, err := r.authenticator()
	if err != nil {
		return nil, err
	}
	image, err := remote.Image(imageRef, remote.WithContext(ctx), remote.WithAuth(auth))
	if err != nil {
		return nil, fmt.Errorf("pull base image %s: %w", reference, err)
	}
	return image, nil
}

func (r *RemoteRegistry) Push(ctx context.Context, reference string, image v1.Image) error {
	imageRef, err := name.ParseReference(reference, name.WeakValidation)
	if err != nil {
		return fmt.Errorf("parse target image reference %q: %w", reference, err)
	}
	if _, ok := imageRef.(name.Tag); !ok {
		return fmt.Errorf("target image reference %q must include a tag", reference)
	}
	auth, err := r.authenticator()
	if err != nil {
		return err
	}
	if err := remote.Write(imageRef, image, remote.WithContext(ctx), remote.WithAuth(auth)); err != nil {
		return fmt.Errorf("push image %s: %w", reference, err)
	}
	return nil
}

func (r *RemoteRegistry) authenticator() (authn.Authenticator, error) {
	if r.username == "" && r.password == "" {
		return authn.Anonymous, nil
	}
	if r.username == "" || r.password == "" {
		return nil, fmt.Errorf("set both BITPACK_REGISTRY_USERNAME and BITPACK_REGISTRY_PASSWORD, or leave both unset for public access")
	}
	return &authn.Basic{Username: r.username, Password: r.password}, nil
}

func (b Builder) Build(ctx context.Context, recipe Recipe) (string, error) {
	if b.Registry == nil {
		return "", fmt.Errorf("image registry client is required")
	}
	if recipe.BaseImage == "" || recipe.Tag == "" || recipe.Workspace == "" {
		return "", fmt.Errorf("base image, target tag, and project workspace are required")
	}
	if b.StoreDirectory == "" {
		return "", fmt.Errorf("local OCI layout directory is required")
	}
	if _, err := name.ParseReference(recipe.Tag, name.WeakValidation); err != nil {
		return "", fmt.Errorf("invalid target image reference %q: %w", recipe.Tag, err)
	}
	workingDir := recipe.WorkingDir
	if workingDir == "" {
		workingDir = "/app"
	}
	if !strings.HasPrefix(workingDir, "/") || path.Clean(workingDir) != workingDir {
		return "", fmt.Errorf("image working directory must be a clean absolute container path")
	}
	base, err := b.Registry.Pull(ctx, recipe.BaseImage)
	if err != nil {
		return "", err
	}
	entries, err := packaging.CollectDirectory(recipe.Workspace, "")
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", fmt.Errorf("project workspace contains no image files")
	}
	for index := range entries {
		entries[index].Path = path.Join(strings.TrimPrefix(workingDir, "/"), entries[index].Path)
	}
	var tarData bytes.Buffer
	if err := packaging.WriteTarEntries(&tarData, entries); err != nil {
		return "", fmt.Errorf("create OCI image layer: %w", err)
	}
	layer, err := newProjectLayer(tarData.Bytes())
	if err != nil {
		return "", fmt.Errorf("prepare OCI project layer: %w", err)
	}
	image, err := mutate.AppendLayers(base, layer)
	if err != nil {
		return "", fmt.Errorf("append project layer to base image: %w", err)
	}
	config, err := image.ConfigFile()
	if err != nil {
		return "", fmt.Errorf("read base image config: %w", err)
	}
	config = config.DeepCopy()
	config.Config.WorkingDir = workingDir
	config.Config.Env = mergeEnvironment(config.Config.Env, recipe.Environment)
	if len(recipe.Entrypoint) > 0 {
		config.Config.Entrypoint = append([]string(nil), recipe.Entrypoint...)
		config.Config.Cmd = append([]string(nil), recipe.Command...)
	} else if len(recipe.Command) > 0 {
		config.Config.Entrypoint = []string{"/bin/sh", "-c"}
		config.Config.Cmd = []string{strings.Join(recipe.Command, " && ")}
	}
	config.Config.Labels = cloneLabels(config.Config.Labels)
	if recipe.Name != "" {
		config.Config.Labels["io.bitpack.name"] = recipe.Name
	}
	image, err = mutate.ConfigFile(image, config)
	if err != nil {
		return "", fmt.Errorf("apply Bitpack image configuration: %w", err)
	}
	image = mutate.MediaType(image, types.OCIManifestSchema1)
	return b.saveLayout(recipe.Tag, image)
}

func (b Builder) Push(ctx context.Context, tag string) error {
	if b.Registry == nil {
		return fmt.Errorf("image registry client is required")
	}
	recordPath := filepath.Join(b.StoreDirectory, referenceKey(tag), "reference.json")
	data, err := os.ReadFile(recordPath)
	if err != nil {
		return fmt.Errorf("read local image reference %s: %w", tag, err)
	}
	var record referenceRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return fmt.Errorf("decode local image reference: %w", err)
	}
	if record.Tag != tag || filepath.Base(record.Layout) != record.Layout {
		return fmt.Errorf("local image reference metadata is invalid for %s", tag)
	}
	imageLayout, err := layout.FromPath(filepath.Join(filepath.Dir(recordPath), record.Layout))
	if err != nil {
		return fmt.Errorf("open local OCI image layout: %w", err)
	}
	index, err := imageLayout.ImageIndex()
	if err != nil {
		return fmt.Errorf("read local OCI image index: %w", err)
	}
	manifest, err := index.IndexManifest()
	if err != nil {
		return fmt.Errorf("read local OCI image descriptors: %w", err)
	}
	for _, descriptor := range manifest.Manifests {
		if descriptor.Annotations[imageRefAnnotation] != tag {
			continue
		}
		image, err := imageLayout.Image(descriptor.Digest)
		if err != nil {
			return fmt.Errorf("read local OCI image: %w", err)
		}
		if err := b.Registry.Push(ctx, tag, image); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("local OCI layout does not contain image tag %s", tag)
}

func (b Builder) saveLayout(tag string, image v1.Image) (string, error) {
	referenceDirectory := filepath.Join(b.StoreDirectory, referenceKey(tag))
	if err := os.MkdirAll(referenceDirectory, 0o700); err != nil {
		return "", fmt.Errorf("create local image directory: %w", err)
	}
	digest, err := image.Digest()
	if err != nil {
		return "", fmt.Errorf("calculate OCI image digest: %w", err)
	}
	layoutName := digest.Hex
	layoutDirectory := filepath.Join(referenceDirectory, layoutName)
	if _, err := os.Stat(layoutDirectory); os.IsNotExist(err) {
		temporary, err := os.MkdirTemp(referenceDirectory, ".oci-build-")
		if err != nil {
			return "", fmt.Errorf("create temporary OCI layout: %w", err)
		}
		defer os.RemoveAll(temporary)
		layoutPath := filepath.Join(temporary, "layout")
		ociLayout, err := layout.Write(layoutPath, empty.Index)
		if err != nil {
			return "", fmt.Errorf("initialize OCI layout: %w", err)
		}
		if err := ociLayout.AppendImage(image, layout.WithAnnotations(map[string]string{
			imageRefAnnotation: tag,
		})); err != nil {
			return "", fmt.Errorf("write image to OCI layout: %w", err)
		}
		if err := os.Rename(layoutPath, layoutDirectory); err != nil {
			return "", fmt.Errorf("save OCI layout: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("inspect existing OCI layout: %w", err)
	}
	record := referenceRecord{Tag: tag, Layout: layoutName}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode local image reference: %w", err)
	}
	if err := os.WriteFile(filepath.Join(referenceDirectory, "reference.json"), data, 0o600); err != nil {
		return "", fmt.Errorf("save local image reference: %w", err)
	}
	return digest.String(), nil
}

func mergeEnvironment(base []string, overrides map[string]string) []string {
	values := make(map[string]string, len(base)+len(overrides))
	for _, entry := range base {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			values[key] = value
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

func cloneLabels(source map[string]string) map[string]string {
	labels := make(map[string]string, len(source)+1)
	for key, value := range source {
		labels[key] = value
	}
	return labels
}

func referenceKey(reference string) string {
	sum := sha256.Sum256([]byte(reference))
	return hex.EncodeToString(sum[:])
}
