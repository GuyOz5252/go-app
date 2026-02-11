package services

import (
	"context"
	"time"

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

func (s *ChatService) ListByUserId(ctx context.Context, userId string) ([]*core.ChatDto, error) {
	return s.chatRepository.ListByUserId(ctx, userId)
}

func (s *ChatService) Create(ctx context.Context, name string, chatMemberIds []string, imageUrl string) (string, error) {
	if (len(chatMemberIds) <= 1) {
		return "", core.ErrMustHaveMoreThanOneMember
	}
	chat := &core.Chat{
		Name:          name,
		ChatMemberIds: chatMemberIds,
		ImageUrl:      imageUrl,
		CreatedAt:     time.Now().UTC(),
	}
	return s.chatRepository.Create(ctx, chat)
}

func (s *ChatService) AddMember(ctx context.Context, chatId, userId string) error {
	ok, err := s.chatRepository.IsMemberInChat(ctx, chatId, userId)
	if err != nil {
		return err
	}
	if ok {
		return core.ErrUserIsAlreadyInChat
	}
	return s.chatRepository.AddMember(ctx, chatId, userId)
}
