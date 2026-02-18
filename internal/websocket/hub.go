package websocket

import (
	"context"
	"sync"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/GuyOz5252/go-app/internal/services"
)

type Hub struct {
	chatService     *services.ChatService
	presenceService *services.PresenceService
	mutex           sync.RWMutex
	clients         map[string]map[*Client]bool
	Register        chan *Client
	Unregister      chan *Client
}

func NewHub(chatService *services.ChatService, presenceService *services.PresenceService) *Hub {
	return &Hub{
		chatService:     chatService,
		presenceService: presenceService,
		clients:         make(map[string]map[*Client]bool),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
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

	isFirstConnection := len(h.clients[client.userId]) == 0

	if _, ok := h.clients[client.userId]; !ok {
		h.clients[client.userId] = make(map[*Client]bool)
	}
	h.clients[client.userId][client] = true

	if isFirstConnection {
		go h.presenceService.SetOnline(context.Background(), client.userId)
	}
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
				go h.presenceService.SetOffline(context.Background(), client.userId)
			}
		}
	}
}

func (h *Hub) sendWSMessage(userId string, wsMessage *core.WSMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, ok := h.clients[userId]; ok {
		for client := range clients {
			select {
			case client.send <- wsMessage:
			default:
				go func() { h.Unregister <- client }()
			}
		}
	}
}

func (h *Hub) broadcastWSMessage(userIds []string, wsMessage *core.WSMessage) {
	for _, userId := range userIds {
		h.sendWSMessage(userId, wsMessage)
	}
}

func (h *Hub) sendWSError(chatId string, userId string, err error) {
	wsErrorMessage := &core.WSMessage{
		Type:    core.ServerError,
		ChatId:  chatId,
		UserId:  userId,
		Payload: err.Error(),
	}
	h.sendWSMessage(userId, wsErrorMessage)
}

func (h *Hub) deliverMessage(ctx context.Context, wsMessage core.WSMessage) {
	wsMessagePayload := wsMessage.Payload.(struct {
		Content   string `json:"content"`
		MediaUrl  string `json:"media_url,omitempty"`
		ReplyToId string `json:"reply_to_id"`
	})

	chatMessage, err := h.chatService.SendMessage(ctx, wsMessage.UserId, wsMessage.ChatId, wsMessagePayload.Content, wsMessagePayload.MediaUrl, wsMessagePayload.ReplyToId)
	if err != nil {
		h.sendWSError(wsMessage.ChatId, wsMessage.UserId, err)
		return
	}

	wsServerAckMessage := &core.WSMessage{
		Type:   core.MessageServerAck,
		ChatId: wsMessage.ChatId,
		UserId: wsMessage.UserId,
		Payload: struct {
			MessageId string `json:"message_id"`
		}{
			MessageId: chatMessage.Id,
		},
	}
	h.sendWSMessage(wsMessage.UserId, wsServerAckMessage)

	chatMembers, err := h.chatService.GetMembers(ctx, wsMessage.ChatId)
	if err != nil {
		h.sendWSError(wsMessage.ChatId, wsMessage.UserId, err)
		return
	}

	wsOutgoingMessage := &core.WSMessage{
		Type:    core.NewMessage,
		UserId:  wsMessage.UserId,
		ChatId:  wsMessage.ChatId,
		Payload: chatMessage,
	}
	h.broadcastWSMessage(chatMembers, wsOutgoingMessage)
}

func (h *Hub) deliverUserAcks(wsMessage core.WSMessage) {
	h.sendWSMessage(wsMessage.UserId, &wsMessage)
}

func (h *Hub) deliverTyping(ctx context.Context, wsMessage core.WSMessage) {
	chatMembers, err := h.chatService.GetMembers(ctx, wsMessage.ChatId)
	if err != nil {
		h.sendWSError(wsMessage.ChatId, wsMessage.UserId, err)
		return
	}

	h.broadcastWSMessage(chatMembers, &wsMessage)
}
