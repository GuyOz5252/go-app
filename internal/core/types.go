package core

import "errors"

var ErrNotFound = errors.New("not found")

type UserRepository interface {
	GetById(id int) (*User, error)
	Create(user *User) (int, error)
	Update(user *User) error
	Delete(id int) error
}
