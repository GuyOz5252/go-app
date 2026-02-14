package core

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")
var ErrUnautherized = errors.New("unautherized")
var ErrQueryNotConfigured = errors.New("query not configured")
var ErrUsernameConflict = errors.New("username already exists")
var ErrEmailConflict = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrMustHaveMoreThanOneMember = errors.New("chat must have more than one member")
var ErrUserIsAlreadyInChat = errors.New("user is already in chat")

type WSMessage struct {
	Type    WSMessageType `json:"message_type"`
	ChatId  string        `json:"chat_id"`
	UserId  string        `json:"user_id"`
	Payload any           `json:"payload,omitempty"`
}

type WSMessageType int

const (
	NewMessage WSMessageType = iota
	MessageServerAck
	MessageUserAck
	MessageUserReadAck
	UserTypingStart
	UserTypingEnd
)

type Cache interface {
	SetKey(ctx context.Context, key string, ttl time.Duration) error
	DeleteKey(ctx context.Context, key string) error
	KeyExists(ctx context.Context, key string) (bool, error)
}

type UserRepository interface {
	GetById(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) (string, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type ChatRepository interface {
	GetById(ctx context.Context, id string) (*Chat, error)
	ListByUserId(ctx context.Context, userId string) ([]*ChatDto, error)
	Create(ctx context.Context, chat *Chat) (string, error)
	AddMember(ctx context.Context, chatId, userId string) error
	RemoveMember(ctx context.Context, chatId, userId string) error
	IsMemberInChat(ctx context.Context, chatId, userId string) (bool, error)
	CreateMessage(ctx context.Context, message *ChatMessage) error
	GetMessages(ctx context.Context, chatId string, limit, offset int) ([]*ChatMessage, error)
}
