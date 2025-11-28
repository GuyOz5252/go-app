package services

import (
	"github.com/GuyOz5252/go-app/internal/core"
)

type UserService struct {
	userRepository core.UserRepository
}

func NewUserService(userRepository core.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetById(id int) (*core.User, error) {
	return s.userRepository.GetById(id)
}

func (s *UserService) Create(user *core.User) (int, error) {
	return s.userRepository.Create(user)
}
