package handlers

import (
	"net/http"

	"github.com/GuyOz5252/go-app/internal/core"
	ws "github.com/GuyOz5252/go-app/internal/services/websocket"
	results "github.com/GuyOz5252/go-app/pkg/api"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub *ws.Hub
	handlerResolver map[string]func(m *core.WSMessage)
}

func NewWebSocketHandler(hub *ws.Hub, hr map[string]func(m *core.WSMessage)) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
		handlerResolver: hr,
	}
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
		results.ApiError(w, r, http.StatusUnauthorized, "unauthorized", "no claims provided")
		return
	}

	userId, ok := claims["userId"].(string)
	if !ok || userId == "" {
		results.ApiError(w, r, http.StatusUnauthorized, "unauthorized", "no userId provided")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		results.ApiError(w, r, http.StatusBadRequest, "failed to upgrade connection", err.Error())
		return
	}

	client := ws.NewClient(userId, h.hub, conn, h.handlerResolver)
	h.hub.Register <- client

	go client.WriteMessages()
	go client.ReadMessages()
}
