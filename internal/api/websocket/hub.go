package websocket

import "sync"

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

func NewHub() *Hub { return &Hub{clients: make(map[string]*Client)} }

func (h *Hub) Add(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.ID] = client
}

func (h *Hub) Remove(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if client, ok := h.clients[id]; ok {
		delete(h.clients, id)
		client.Close()
	}
}

func (h *Hub) Broadcast(event []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.clients {
		select {
		case client.Events <- event:
		default:
		}
	}
}
