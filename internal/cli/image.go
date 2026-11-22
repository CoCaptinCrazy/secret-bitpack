package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	ociimage "github.com/CoCaptinCrazy/secret-bitpack/internal/builder/image"
	devruntime "github.com/CoCaptinCrazy/secret-bitpack/internal/dev/runtime"
)

func bitpackImage(output io.Writer, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: image <build|push> [arguments]")
	}
	switch args[0] {
	case "build":
		return bitpackImageBuild(output, args[1:])
	case "push":
		if len(args) != 2 {
			return fmt.Errorf("usage: image push <registry/repository:tag>")
		}
		builder := ociimage.Builder{
			Registry:       ociimage.NewRemoteRegistry(),
			StoreDirectory: dataPath("images"),
		}
		if err := builder.Push(context.Background(), args[1]); err != nil {
			return err
		}
		fmt.Fprintf(output, "Image pushed: %s\n", args[1])
		return nil
	default:
		return fmt.Errorf("unknown image action %q; use build or push", args[0])
	}
}

func bitpackImageBuild(output io.Writer, args []string) error {
	source := "Bitfile"
	tag := ""
	index := 0
	if len(args) > 0 && !strings.HasPrefix(args[0], "--") {
		source = args[0]
		index++
	}
	for index < len(args) {
		switch args[index] {
		case "--tag":
			if index+1 == len(args) {
				return fmt.Errorf("--tag requires a value")
			}
			index++
			tag = args[index]
		default:
			return fmt.Errorf("unknown image build option %q", args[index])
		}
		index++
	}
	if tag == "" {
		return fmt.Errorf("usage: image build [Bitfile] --tag <registry/repository:tag>")
	}
	blueprint, workspace, _, err := loadBlueprintInput(source)
	if err != nil {
		return err
	}
	baseImage := blueprint.Directives["base"]
	if baseImage == "" {
		return fmt.Errorf("Bitfile requires a base directive, for example: base \"node:22-alpine\"")
	}
	if len(blueprint.Assemble) > 0 {
		environment, err := blueprintEnvironment(blueprint)
		if err != nil {
			return err
		}
		if err := devruntime.RunCommands(context.Background(), blueprint.Assemble, environment,
			workspace, os.Stdin, output, os.Stderr); err != nil {
			return fmt.Errorf("run image build commands: %w", err)
		}
	}
	var entrypoint []string
	if raw := blueprint.Directives["entrypoint"]; raw != "" {
		if err := json.Unmarshal([]byte(raw), &entrypoint); err != nil || len(entrypoint) == 0 {
			return fmt.Errorf("Bitfile entrypoint must be a non-empty JSON string array")
		}
	}
	builder := ociimage.Builder{
		Registry:       ociimage.NewRemoteRegistry(),
		StoreDirectory: dataPath("images"),
	}
	digest, err := builder.Build(context.Background(), ociimage.Recipe{
		BaseImage: baseImage,
		Tag:       tag,
		Name:      blueprint.Name,
		Workspace: workspace,
		WorkingDir: func() string {
			if value := blueprint.Directives["workdir"]; value != "" {
				return value
			}
			return "/app"
		}(),
		Environment: blueprint.Environment,
		Entrypoint:  entrypoint,
		Command:     blueprint.Ignite,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(output, "OCI image built: %s\nDigest: %s\n", tag, digest)
	fmt.Fprintf(output, "Stored locally in Bitpack's OCI layout; push with: bitpack image push %s\n", tag)
	return nil
}
