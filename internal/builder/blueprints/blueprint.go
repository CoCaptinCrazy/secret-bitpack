package blueprints

type Blueprint struct {
	Name        string
	Version     string
	Using       string
	Context     string
	Layers      []Layer
	Metadata    map[string]any
	Environment map[string]string
	Secrets     []string
	Health      []string
	Assemble    []string
	Verify      []string
	Ignite      []string
	Sections    map[string][]string
	Directives  map[string]string
}

type Layer struct {
	Name    string
	Content string
}

func New(name, version string) *Blueprint {
	return &Blueprint{
		Name:        name,
		Version:     version,
		Metadata:    make(map[string]any),
		Environment: make(map[string]string),
		Secrets:     []string{},
		Health:      []string{},
		Assemble:    []string{},
		Verify:      []string{},
		Ignite:      []string{},
		Sections:    make(map[string][]string),
		Directives:  make(map[string]string),
	}
}

func (b *Blueprint) AddLayer(layer Layer) {
	b.Layers = append(b.Layers, layer)
}
