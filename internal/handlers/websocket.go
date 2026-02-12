package handlers

import (
	"net/http"

	ws "github.com/GuyOz5252/go-app/internal/services/websocket"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	if claims == nil {
		// TODO: Log
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId, ok := claims["userId"].(string)
	if !ok || userId == "" {
		// TODO: Log
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO: Log
		return
	}

	client := ws.NewClient(h.hub, conn, userId)
	h.hub.Register(client)

	go client.WriteMessages()
	go client.ReadMessages()
}
