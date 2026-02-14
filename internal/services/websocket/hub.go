package websocket

import (
	"sync"

	"github.com/GuyOz5252/go-app/internal/core"
)

type Hub struct {
	mutex      sync.RWMutex
	clients    map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client.userId]; !ok {
		h.clients[client.userId] = make(map[*Client]bool)
	}
	h.clients[client.userId][client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if userClients, ok := h.clients[client.userId]; ok {
		if _, ok := userClients[client]; ok {
			delete(userClients, client)
			close(client.send)
			client.connection.Close()
			if len(userClients) == 0 {
				delete(h.clients, client.userId)
			}
		}
	}
}

func (h *Hub) PublishMessage(userId string, message *core.WSMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, ok := h.clients[userId]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				go func() { h.Unregister <- client }()
			}
		}
	}
}
