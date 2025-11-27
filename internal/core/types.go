package core

import "errors"

var NotFoundError = errors.New("not found")

type UserRepository interface {
	GetById(id int) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
}
