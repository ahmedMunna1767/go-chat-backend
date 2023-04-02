package websocket

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Pool struct {
	mutex   sync.RWMutex
	Clients map[uuid.UUID]*Client
}

func NewPool() *Pool {
	return &Pool{
		mutex:   sync.RWMutex{},
		Clients: make(map[uuid.UUID]*Client),
	}
}

func (p *Pool) AddConnection(client *Client) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Clients[client.ID] = client
}

func (p *Pool) RemoveConnection(client *Client) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	delete(p.Clients, client.ID)
}

func (p *Pool) GetClientByID(id uuid.UUID) (*Client, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	client, ok := p.Clients[id]
	if !ok {
		return nil, fmt.Errorf("client not found")
	}
	return client, nil
}
