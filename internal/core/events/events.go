package events

import "sync"

type Handler func(Event)

type Event struct {
	Type    string
	Payload any
}

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func New() *Bus {
	return &Bus{handlers: make(map[string][]Handler)}
}

func (b *Bus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *Bus) Publish(event Event) {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[event.Type]...)
	b.mu.RUnlock()
	for _, handler := range handlers {
		handler(event)
	}
}
