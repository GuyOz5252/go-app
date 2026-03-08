package websocket

import (
	"context"
	"sync"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/GuyOz5252/go-app/internal/services"
)

type Hub struct {
	chatService            *services.ChatService
	presenceService        *services.PresenceService
	userConnectionsService *services.UserConnectionsService
	mutex                  sync.RWMutex
	clients                map[string]map[*Client]bool
	Register               chan *Client
	Unregister             chan *Client
}

func NewHub(chatService *services.ChatService, presenceService *services.PresenceService, userConnectionsService *services.UserConnectionsService) *Hub {
	return &Hub{
		chatService:            chatService,
		presenceService:        presenceService,
		userConnectionsService: userConnectionsService,
		clients:                make(map[string]map[*Client]bool),
		Register:               make(chan *Client),
		Unregister:             make(chan *Client),
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
		go h.presenceService.SetOnline(context.Background(), client.userId)
	}
	h.clients[client.userId][client] = true
	h.userConnectionsService.AddConnection(context.Background(), client.userId, client.connection.RemoteAddr().String())
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
				h.userConnectionsService.DeleteConnection(context.Background(), client.userId, client.connection.RemoteAddr().String())
				go h.presenceService.SetOffline(context.Background(), client.userId)
			}
		}
	}
}

func (h *Hub) sendWSMessage(userId string, wsMessage *core.WSMessage) {
	h.mutex.RLock()
	if clients, ok := h.clients[userId]; ok {
		for client := range clients {
			select {
			case client.send <- wsMessage:
			default:
				go func() { h.Unregister <- client }()
			}
		}
	}
	h.mutex.RUnlock()

	if connections, err := h.userConnectionsService.GetUserConnections(context.Background(), userId); err != nil {
		for _, connectionId := range connections {
			// TODO: check if connection id is of this server
			// TODO: forward to connected server using redis pub/sub
			println("Connection ID:", connectionId)
		}
	}
}

func (h *Hub) broadcastWSMessage(userIds []string, wsMessage *core.WSMessage) {
	for _, userId := range userIds {
		h.sendWSMessage(userId, wsMessage)
	}
}

func (h *Hub) sendWSError(userId string, err error) {
	wsErrorMessage := &core.WSMessage{
		Type:              core.ServerError,
		DestinationUserId: userId,
		Payload:           err.Error(),
	}
	h.sendWSMessage(userId, wsErrorMessage)
}

func (h *Hub) sendWSMessageToDestinationChat(ctx context.Context, wsMessage *core.WSMessage) {
	chatMembers, err := h.chatService.GetMembers(ctx, wsMessage.DestinationChatId)
	if err != nil {
		h.sendWSError(wsMessage.InitiatorUserId, err)
		return
	}

	h.broadcastWSMessage(chatMembers, wsMessage)
}

func (h *Hub) deliverMessage(ctx context.Context, wsMessage *core.WSMessage) {
	wsMessagePayload := wsMessage.Payload.(struct {
		Content   string `json:"content"`
		MediaUrl  string `json:"media_url,omitempty"`
		ReplyToId string `json:"reply_to_id"`
	})

	chatMessage, err := h.chatService.SendMessage(ctx, wsMessage.InitiatorUserId, wsMessage.DestinationChatId, wsMessagePayload.Content, wsMessagePayload.MediaUrl, wsMessagePayload.ReplyToId)
	if err != nil {
		h.sendWSError(wsMessage.InitiatorUserId, err)
		return
	}

	wsServerAckMessage := &core.WSMessage{
		Type:              core.MessageServerAck,
		DestinationUserId: wsMessage.InitiatorUserId,
		Payload: struct {
			MessageId string `json:"message_id"`
		}{
			MessageId: chatMessage.Id,
		},
	}
	h.sendWSMessage(wsMessage.InitiatorUserId, wsServerAckMessage)

	wsOutgoingMessage := &core.WSMessage{
		Type:              core.NewMessage,
		InitiatorUserId:   wsMessage.InitiatorUserId,
		DestinationChatId: wsMessage.DestinationChatId,
		Payload:           chatMessage,
	}
	h.sendWSMessageToDestinationChat(ctx, wsOutgoingMessage)
}
