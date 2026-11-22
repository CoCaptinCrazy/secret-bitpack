package scheduler

func Place(nodes []Node, requirement Requirement) (Node, bool) {
	for _, node := range nodes {
		if Matches(node, requirement) {
			return node, true
		}
	}
	return Node{}, false
}
