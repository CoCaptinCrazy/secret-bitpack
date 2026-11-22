package websocket

import "sync"

type Client struct {
	ID     string
	Events chan []byte
	once   sync.Once
}

func NewClient(id string, buffer int) *Client {
	return &Client{ID: id, Events: make(chan []byte, buffer)}
}

func (c *Client) Close() {
	c.once.Do(func() { close(c.Events) })
}
