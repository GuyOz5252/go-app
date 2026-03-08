package services

import (
	"context"
	"fmt"

	"github.com/GuyOz5252/go-app/internal/core"
)

type UserConnectionsService struct {
	cache core.Cache
}

func NewUserConnectionsService(cache core.Cache) *UserConnectionsService {
	return &UserConnectionsService{
		cache: cache,
	}
}

func (s *UserConnectionsService) AddConnection(ctx context.Context, userId, connectionId string) error {
	key := fmt.Sprintf("connection:%s", userId)
	return s.cache.AddToSet(ctx, key, connectionId)
}

func (s *UserConnectionsService) DeleteConnection(ctx context.Context, userId, connectionId string) error {
	key := fmt.Sprintf("connection:%s", userId)
	return s.cache.RemoveFromSet(ctx, key, connectionId)
}

func (s *UserConnectionsService) GetUserConnections(ctx context.Context, userId string) ([]string, error) {
	key := fmt.Sprintf("connection:%s", userId)
	return s.cache.GetSet(ctx, key)
}
