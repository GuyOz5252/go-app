package services

import (
	"context"

	"github.com/GuyOz5252/go-app/internal/core"
	auth "github.com/GuyOz5252/go-app/pkg"
)

type UserService struct {
	userRepository core.UserRepository
}

func NewUserService(userRepository core.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetById(ctx context.Context, id string) (*core.User, error) {
	return s.userRepository.GetById(ctx, id)
}

func (s *UserService) Register(ctx context.Context, username, email, rawPassword string) (string, error) {
	usernameExists, err := s.userRepository.ExistsByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if usernameExists {
		return "", core.ErrUsernameConflict
	}

	emailExists, err := s.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if emailExists {
		return "", core.ErrEmailConflict
	}

	hash, err := auth.HashPassword(rawPassword)
	if err != nil {
		return "", err
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

	if !auth.CheckPassword(passwordStr, user.PasswordHash) {
		return nil, core.ErrInvalidCredentials
	}

	return user, nil
}
