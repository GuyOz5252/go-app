package core

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrQueryNotConfigured = errors.New("query not configured")
var ErrUsernameConflict = errors.New("username already exists")
var ErrEmailConflict = errors.New("email already exists")

type UserRepository interface {
	GetById(ctx context.Context, id int) (*User, error)
	Create(ctx context.Context, user *User) (int, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
