package parser

import "fmt"

func (g *Graph) Validate() error {
	for step, dependencies := range g.Steps {
		for _, dependency := range dependencies {
			if _, ok := g.Steps[dependency]; !ok {
				return fmt.Errorf("step %q depends on missing step %q", step, dependency)
			}
		}
	}
	return nil
}
