package services

import (
	"context"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/GuyOz5252/go-app/pkg/password"
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

func (s *UserService) Create(ctx context.Context, username, email, rawPassword string) (int, error) {
	usernameExists, err := s.userRepository.ExistsByUsername(ctx, username)
	if err != nil {
		return -1, err
	}
	if usernameExists {
		return -1, core.ErrUsernameConflict
	}

	emailExists, err := s.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return -1, err
	}
	if emailExists {
		return -1, core.ErrEmailConflict
	}

	hash, err := password.HashPassword(rawPassword)
	if err != nil {
		return -1, err
	}

	user := &core.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
	}

	return s.userRepository.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, email, passwordStr string) (*core.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, core.ErrInvalidCredentials
		}
		return nil, err
	}

	if !password.CheckPassword(passwordStr, user.PasswordHash) {
		return nil, core.ErrInvalidCredentials
	}

	return user, nil
}
