package mesh

import "sync"

type Peer struct {
	ID      string
	Address string
}

type Mesh struct {
	mu    sync.RWMutex
	peers map[string]Peer
}

func New() *Mesh {
	return &Mesh{peers: make(map[string]Peer)}
}

func (m *Mesh) Join(peer Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.peers[peer.ID] = peer
}

func (m *Mesh) Peers() []Peer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]Peer, 0, len(m.peers))
	for _, peer := range m.peers {
		result = append(result, peer)
	}
	return result
}
