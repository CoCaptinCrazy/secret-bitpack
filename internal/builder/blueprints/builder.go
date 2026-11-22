package blueprints

import "fmt"

type Builder struct{}

func (b *Builder) Build(name, version string, layers []Layer) (*Blueprint, error) {
	if name == "" {
		return nil, fmt.Errorf("blueprint name is required")
	}
	blueprint := New(name, version)
	for _, layer := range layers {
		blueprint.AddLayer(layer)
	}
	return blueprint, nil
}
