package sessionstorage

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisSessionStorage struct {
	config      *models.ConfigOptions
	redisClient *redis.Client
}

type redisStorageElement map[string]any

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
	expTime := s.config.GetAccessJWTExpiration()

	sessionData := make(redisStorageElement)

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

	return s.redisClient.Expire(ctx, userID.String(), s.config.GetAccessJWTExpiration()).Err()
}

func (s *RedisSessionStorage) CheckSession(ctx context.Context, userID uuid.UUID) error {
	exists, err := s.redisClient.Exists(ctx, userID.String()).Result()
	if err != nil {
		return err
	}
	if exists == 0 { // if session never existed, or reached ttl
		return DoesNotExists
	}

	return nil
}
