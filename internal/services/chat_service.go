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

	s.sendMessage(chatId, userId, chatMessage)

	return chatMessage, nil
}

func (s *ChatService) SendMessageWsHandler(wsMessage *core.WSMessage) {
	isMember, err := s.chatRepository.IsMemberInChat(context.Background(), wsMessage.ChatId, wsMessage.UserId)
	if err != nil || !isMember {
		return
	}

	wsMessagePayload := wsMessage.Payload.(struct {
		Content   string `json:"content"`
		MediaUrl  string `json:"media_url,omitempty"`
		ReplyToId string `json:"reply_to_id"`
	})

	chatMessage := core.ChatMessage{
		UserId:    wsMessage.UserId,
		ChatId:    wsMessage.ChatId,
		Content:   wsMessagePayload.Content,
		MediaUrl:  wsMessagePayload.MediaUrl,
		ReplyToId: wsMessagePayload.ReplyToId,
		CreatedAt: time.Now().UTC(),
	}

	s.sendMessage(wsMessage.ChatId, wsMessage.UserId, &chatMessage)
}

func (s *ChatService) sendMessage(chatId string, userId string, chatMessage *core.ChatMessage) {
	err := s.chatRepository.CreateMessage(context.Background(), chatMessage)
	if err != nil {
		return
	}

	s.hub.PublishMessage(userId, &core.WSMessage{
		Type: core.MessageServerAck,
		Payload: struct {
			ChatId    string `json:"chat_id"`
			MessageId string `json:"message_id"`
		}{
			ChatId:    chatId,
			MessageId: chatMessage.Id,
		},
	})

	newWsMessage := &core.WSMessage{
		Type:    core.NewMessage,
		ChatId:  chatId,
		UserId:  userId,
		Payload: chatMessage,
	}
	s.PublishWSMessageToChat(newWsMessage)
}

func (s *ChatService) PublishWSMessageToChat(wsMessage *core.WSMessage) {
	chat, err := s.chatRepository.GetById(context.Background(), wsMessage.ChatId)
	if err != nil {
		return
	}

	for _, userId := range chat.ChatMemberIds {
		if userId != wsMessage.UserId {
			s.hub.PublishMessage(userId, wsMessage)
		}
	}
}
