package services

import (
	"context"
	"fmt"
	"time"

	"github.com/GuyOz5252/go-app/internal/core"
)

type PresenceService struct {
	cache core.Cache
}

func NewPresenceService(cache core.Cache) *PresenceService {
	return &PresenceService{
		cache: cache,
	}
}

func (s *PresenceService) SetOnline(ctx context.Context, userId string) error {
	key := fmt.Sprintf("presence:online:%s", userId)
	return s.cache.SetKey(ctx, key, time.Minute*5)
}

func (s *PresenceService) SetOffline(ctx context.Context, userId string) error {
	key := fmt.Sprintf("presence:online:%s", userId)
	return s.cache.DeleteKey(ctx, key)
}

func (s *PresenceService) IsOnline(ctx context.Context, userId string) (bool, error) {
	key := fmt.Sprintf("presence:online:%s", userId)
	return s.cache.KeyExists(ctx, key)
}
