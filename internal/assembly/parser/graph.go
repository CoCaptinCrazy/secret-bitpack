package parser

type Graph struct {
	Steps map[string][]string
}

func NewGraph() *Graph { return &Graph{Steps: make(map[string][]string)} }

func (g *Graph) AddStep(name string, dependencies ...string) {
	g.Steps[name] = append([]string(nil), dependencies...)
}
