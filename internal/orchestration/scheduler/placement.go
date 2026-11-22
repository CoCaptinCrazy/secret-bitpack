package scheduler

type Requirement struct {
	Labels map[string]string
}

type Node struct {
	ID     string
	Labels map[string]string
}

func Matches(node Node, requirement Requirement) bool {
	for key, value := range requirement.Labels {
		if node.Labels[key] != value {
			return false
		}
	}
	return true
}
