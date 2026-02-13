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

	chatMessage := &core.ChatMessage{
		UserId:    userId,
		ChatId:    chatId,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	err = s.chatRepository.CreateMessage(ctx, chatMessage)
	if err != nil {
		return nil, err
	}

	wsAckMessage := core.WSMessage{
		Type:    "message:ack",
		Payload: chatMessage.Id,
	}
	s.hub.PublishMessage(chatMessage.UserId, &wsAckMessage)

	chat, err := s.chatRepository.GetById(ctx, chatId)
	if err != nil {
		return nil, err
	}

	wsMessage := core.WSMessage{
		Type:    "message:new",
		Payload: chatMessage,
	}
	for _, userId := range chat.ChatMemberIds {
		if userId != chatMessage.UserId {
			s.hub.PublishMessage(userId, &wsMessage)
		}
	}

	return chatMessage, nil
}

func (s *ChatService) NotifyUserTyping(m *core.WSMessage) {
	typingMessage, ok := m.Payload.(struct {
		UserId string `json:"user_id"`
		ChatId string `json:"chat_id"`
	})
	if !ok {
		return
	}

	chat, err := s.chatRepository.GetById(nil, typingMessage.ChatId)
	if err != nil {
		return
	}

	for _, userId := range chat.ChatMemberIds {
		if userId != typingMessage.UserId {
			s.hub.PublishMessage(userId, m)
		}
	}
}
