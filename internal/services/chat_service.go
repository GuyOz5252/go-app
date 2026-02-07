package services

import (
	"context"

	"github.com/GuyOz5252/go-app/internal/core"
)

type ChatService struct {
	chatRepository core.ChatRepository
}

func NewChatService(chatRepository core.ChatRepository) *ChatService {
	return &ChatService{
		chatRepository: chatRepository,
	}
}

func (s *ChatService) GetById(ctx context.Context, id string) (*core.Chat, error) {
	return s.chatRepository.GetById(ctx, id)
}

func (s *ChatService) ListByUserId(ctx context.Context, userId string) ([]*core.Chat, error) {
	return s.chatRepository.ListByUserId(ctx, userId)
}
