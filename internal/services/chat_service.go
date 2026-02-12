package services

import (
	"context"
	"time"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/GuyOz5252/go-app/internal/services/websocket"
)

type ChatService struct {
	chatRepository core.ChatRepository
	hub            *websocket.Hub
}

func NewChatService(chatRepository core.ChatRepository, hub *websocket.Hub) *ChatService {
	return &ChatService{
		chatRepository: chatRepository,
		hub:            hub,
	}
}

func (s *ChatService) GetById(ctx context.Context, id string) (*core.Chat, error) {
	return s.chatRepository.GetById(ctx, id)
}

func (s *ChatService) ListByUserId(ctx context.Context, userId string) ([]*core.ChatDto, error) {
	return s.chatRepository.ListByUserId(ctx, userId)
}

func (s *ChatService) Create(ctx context.Context, name string, chatMemberIds []string, imageUrl string) (string, error) {
	if len(chatMemberIds) <= 1 {
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

func (s *ChatService) SendMessage(ctx context.Context, userId, chatId, content string) (*core.ChatMessage, error) {
	isMember, err := s.chatRepository.IsMemberInChat(ctx, chatId, userId)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, core.ErrUnautherized
	}

	msg := &core.ChatMessage{
		UserId:    userId,
		ChatId:    chatId,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	err = s.chatRepository.CreateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	chat, err := s.chatRepository.GetById(ctx, chatId)
	if err != nil {
		return nil, err
	}

	wsMessage := core.WSMessage{
		Type:    "message",
		Payload: msg,
	}
	for _, userId := range chat.ChatMemberIds {
		s.hub.PublishMessage(userId, &wsMessage)
	}

	return msg, nil
}
