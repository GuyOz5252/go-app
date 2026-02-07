package core

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrQueryNotConfigured = errors.New("query not configured")
var ErrUsernameConflict = errors.New("username already exists")
var ErrEmailConflict = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

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
}
