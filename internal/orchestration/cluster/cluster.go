package cluster

import "sync"

type Node struct {
	ID     string
	Status string
}

type Cluster struct {
	mu    sync.RWMutex
	nodes map[string]Node
}

func New() *Cluster {
	return &Cluster{nodes: make(map[string]Node)}
}

func (c *Cluster) Add(node Node) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nodes[node.ID] = node
}

func (c *Cluster) Nodes() []Node {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]Node, 0, len(c.nodes))
	for _, node := range c.nodes {
		result = append(result, node)
	}
	return result
}
