package sessionstorage

import (
	"context"
	"encoding"
	"encoding/json"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisSessionStorage struct {
	config      *models.ConfigOptions
	redisClient *redis.Client
}

type redisStorageElement map[string]any

func (se redisStorageElement) MarshalBinary() (data []byte, err error) {
	return json.Marshal(se)
}

func NewRedisSessionStorage(config *models.ConfigOptions) *RedisSessionStorage {
	return &RedisSessionStorage{
		config: config,
		redisClient: redis.NewClient(&redis.Options{
			Addr:     config.GetRedisAddress(),
			Password: config.GetRedisPassword(),
			DB:       config.GetRedisDB(),
		}),
	}
}

func (s *RedisSessionStorage) CreateSession(ctx context.Context, userID uuid.UUID) error {
	expTime := s.config.GetRefreshJWTExpiration()
	var sessionData encoding.BinaryMarshaler

	sessionData = make(redisStorageElement)

	return s.redisClient.Set(ctx, userID.String(), sessionData, expTime).Clone().Err()
}

func (s *RedisSessionStorage) DeleteSession(ctx context.Context, userID uuid.UUID) error {
	err := s.CheckSession(ctx, userID)
	if err != nil {
		return err
	}

	return s.redisClient.Del(ctx, userID.String()).Err()
}

func (s *RedisSessionStorage) ExtendSession(ctx context.Context, userID uuid.UUID) error {
	err := s.CheckSession(ctx, userID)
	if err != nil {
		return domainerrors.InvalidCredentials
	}
	return s.redisClient.Expire(ctx, userID.String(), s.config.GetRefreshJWTExpiration()).Err()
}

func (s *RedisSessionStorage) CheckSession(ctx context.Context, userID uuid.UUID) error {
	exists, err := s.redisClient.Exists(ctx, userID.String()).Result()
	if err != nil {
		return err
	}
	if exists == 0 { // if session never existed, or reached ttl
		return domainerrors.InvalidCredentials
	}

	return nil
}
