package parser

type Plan struct {
	Name     string
	Version  string
	Author   string
	Steps    []string
	Forge    []string
	Verify   []string
	Publish  []string
	Assemble []string
	Ignite   []string
	Parallel bool
	FailFast bool
	Requires []string
	Tags     []string
}

func (p *Plan) Add(step string) { p.Steps = append(p.Steps, step) }
