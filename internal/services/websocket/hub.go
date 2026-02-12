package websocket

import (
	"sync"
)

type Hub struct {
	clients    map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			if _, ok := h.clients[client.userId]; !ok {
				h.clients[client.userId] = make(map[*Client]bool)
			}
			h.clients[client.userId][client] = true
			h.mutex.Unlock()
		case client := <-h.unregister:
			h.mutex.Lock()
			if userClients, ok := h.clients[client.userId]; ok {
				if _, ok := userClients[client]; ok {
					delete(userClients, client)
					client.connection.Close()
					if len(userClients) == 0 {
						delete(h.clients, client.userId)
					}
				}
			}
			h.mutex.Unlock()
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) PublishMessage(userId string, message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if clients, ok := h.clients[userId]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}
