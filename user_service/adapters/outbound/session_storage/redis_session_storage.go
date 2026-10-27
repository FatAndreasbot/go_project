package sessionstorage

import (
	"context"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type RedisSessionStorage struct {
	config *models.ConfigOptions
}

func NewRedisSessionStorage(config *models.ConfigOptions) *RedisSessionStorage {
	return &RedisSessionStorage{
		config: config,
	}
}

func (r *RedisSessionStorage) CreateSession(ctx context.Context, userID uuid.UUID) error {
	return domainerrors.NotImplemented
}

func (r *RedisSessionStorage) DeleteSession(ctx context.Context, userID uuid.UUID) error {
	return domainerrors.NotImplemented
}

func (r *RedisSessionStorage) ExtendSession(ctx context.Context, userID uuid.UUID) error {
	return domainerrors.NotImplemented
}

func (r *RedisSessionStorage) CheckSession(ctx context.Context, userID uuid.UUID) error {
	return domainerrors.NotImplemented
}
