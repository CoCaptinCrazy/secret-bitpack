package parser

type Service struct {
	Name    string
	Command string
}

type Environment struct {
	Name      string
	Version   string
	Runtime   string
	Workspace string
	Services  []Service
	Ports     []int
	Watch     []string
	Variables map[string]string
	Secrets   []string
	Startup   []string
	Shell     []string
}
