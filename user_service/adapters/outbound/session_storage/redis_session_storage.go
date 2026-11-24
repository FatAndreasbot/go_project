package sessionstorage

import (
	"context"
	"encoding"
	"encoding/json"
	"log"

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

func unmarshalStorageElement(data string) (*redisStorageElement, error) {
	var storageElement *redisStorageElement
	err := json.Unmarshal([]byte(data), storageElement)
	if err != nil {
		return nil, err
	}
	return storageElement, nil
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

// TODO CreateSession should return a unique sessionID
func (s *RedisSessionStorage) CreateSession(ctx context.Context) (uuid.UUID, error) {
	newSessionID, err := uuid.NewUUID()
	if err != nil {
		log.Default().Println(err)
		return uuid.Nil, err
	}

	expTime := s.config.GetRefreshJWTExpiration()
	var sessionData encoding.BinaryMarshaler

	sessionData = make(redisStorageElement)

	return newSessionID, s.redisClient.Set(ctx, newSessionID.String(), sessionData, expTime).Clone().Err()
}

func (s *RedisSessionStorage) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	err := s.CheckSession(ctx, sessionID)
	if err != nil {
		return err
	}

	return s.redisClient.Del(ctx, sessionID.String()).Err()
}

func (s *RedisSessionStorage) ExtendSession(ctx context.Context, sessionID uuid.UUID) error {
	err := s.CheckSession(ctx, sessionID)
	if err != nil {
		return domainerrors.InvalidCredentials
	}
	return s.redisClient.Expire(ctx, sessionID.String(), s.config.GetRefreshJWTExpiration()).Err()
}

func (s *RedisSessionStorage) CheckSession(ctx context.Context, sessionID uuid.UUID) error {
	exists, err := s.redisClient.Exists(ctx, sessionID.String()).Result()
	if err != nil {
		return err
	}
	if exists == 0 { // if session never existed, or reached ttl
		return domainerrors.InvalidCredentials
	}

	return nil
}

func (s *RedisSessionStorage) SetData(ctx context.Context, sessionID uuid.UUID, key string, data any) error {
	if err := s.CheckSession(ctx, sessionID); err != nil {
		return err
	}

	sessionDataString, err := s.redisClient.Get(ctx, sessionID.String()).Result()
	if err != nil {
		return err
	}

	sessionData, err := unmarshalStorageElement(sessionDataString)
	if err != nil {
		return err
	}

	(*sessionData)[key] = data
	expTime := s.config.GetRefreshJWTExpiration()

	return s.redisClient.Set(ctx, sessionID.String(), sessionData, expTime).Err()
}

func (s *RedisSessionStorage) GetData(ctx context.Context, sessionKey uuid.UUID, key string) (any, error) {
	return nil, domainerrors.NotImplemented
}
