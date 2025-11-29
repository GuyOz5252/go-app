package core

import "errors"

var ErrNotFound = errors.New("not found")
var ErrQueryNotConfigured = errors.New("query not configured")
var ErrUsernameConflict = errors.New("username already exists")
var ErrEmailConflict = errors.New("email already exists")

type UserRepository interface {
	GetById(id int) (*User, error)
	Create(user *User) (int, error)
	Update(user *User) error
	Delete(id int) error
	ExistsByUsername(email string) (bool, error)
	ExistsByEmail(email string) (bool, error)
}
