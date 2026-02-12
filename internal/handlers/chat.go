package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GuyOz5252/go-app/internal/services"
	results "github.com/GuyOz5252/go-app/pkg/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type ChatHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

type CreateChatRequest struct {
	Name          string   `json:"name"`
	ChatMemberIds []string `json:"chat_member_ids"`
	ImageUrl      string   `json:"image_url"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}

func (h *ChatHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		results.ApiError(w, r, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	id, err := h.chatService.Create(r.Context(), req.Name, req.ChatMemberIds, req.ImageUrl)
	if err != nil {
		results.ApiError(w, r, http.StatusInternalServerError, "failed to create chat", err.Error())
		return
	}

	results.ApiResponse(w, r, http.StatusCreated, map[string]string{"id": id})
}

func (h *ChatHandler) List(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId := claims["userId"].(string)

	chats, err := h.chatService.ListByUserId(r.Context(), userId)
	if err != nil {
		results.ApiError(w, r, http.StatusInternalServerError, "failed to list chats", err.Error())
		return
	}

	results.ApiResponse(w, r, http.StatusOK, chats)
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	chatId := chi.URLParam(r, "chatId")
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId := claims["userId"].(string)

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		results.ApiError(w, r, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	msg, err := h.chatService.SendMessage(r.Context(), userId, chatId, req.Content)
	if err != nil {
		results.ApiError(w, r, http.StatusInternalServerError, "failed to send message", err.Error())
		return
	}

	results.ApiResponse(w, r, http.StatusCreated, msg)
}
