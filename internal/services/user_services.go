package services

import (
	"context"

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

func (s *UserService) GetById(ctx context.Context, id int) (*core.User, error) {
	return s.userRepository.GetById(ctx, id)
}

func (s *UserService) Create(ctx context.Context, user *core.User) (int, error) {
	usernameExists, err := s.userRepository.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return -1, err
	}
	if usernameExists {
		return -1, core.ErrUsernameConflict
	}

	emailExists, err := s.userRepository.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return -1, err
	}
	if emailExists {
		return -1, core.ErrEmailConflict
	}

	return s.userRepository.Create(ctx, user)
}
