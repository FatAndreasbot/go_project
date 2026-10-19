package sessionstorage

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/infra/config"
	"github.com/google/uuid"
)

type MapSessionStorage struct {
	store map[uuid.UUID]storageElement
	mu    sync.Mutex
}

type storageElement struct {
	expiresAt time.Time
	data      map[string]any
}

func NewMapSessionStorate() *MapSessionStorage {
	return &MapSessionStorage{
		store: make(map[uuid.UUID]storageElement),
		mu:    sync.Mutex{},
	}
}

func (s *MapSessionStorage) CreateSession(ctx context.Context, userID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sessionCreatedAt := time.Now()
	expTimeString, err := config.GetConfig(config.AccessJWTExpiration)
	if err != nil {
		return err
	}
	expTimeSeconds, err := strconv.Atoi(expTimeString)
	if err != nil {
		return err
	}
	expTime := time.Second * time.Duration(expTimeSeconds)
	expiresAt := sessionCreatedAt.Add(expTime)

	s.store[userID] = storageElement{
		expiresAt: expiresAt,
		data:      make(map[string]any),
	}

	return nil
}

func (s *MapSessionStorage) DeleteSession(ctx context.Context, userID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.store[userID]
	if !ok {
		return DoesNotExists
	}
	delete(s.store, userID)
	if data.expiresAt.After(time.Now()) {
		return Expired
	}

	return nil
}

func (s *MapSessionStorage) ExtendSession(ctx context.Context, userID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.store[userID]
	if !ok {
		return DoesNotExists
	}
	if data.expiresAt.After(time.Now()) {
		delete(s.store, userID)
		return Expired
	}

	defer s.mu.Unlock()
	sessionCreatedAt := time.Now()
	expTimeString, err := config.GetConfig(config.AccessJWTExpiration)
	if err != nil {
		return err
	}
	expTimeSeconds, err := strconv.Atoi(expTimeString)
	if err != nil {
		return err
	}
	expTime := time.Second * time.Duration(expTimeSeconds)
	newExpiresAt := sessionCreatedAt.Add(expTime)

	data.expiresAt = newExpiresAt
	s.store[userID] = data

	return nil
}

func (s *MapSessionStorage) CheckSession(ctx context.Context, userID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.store[userID]
	if !ok {
		return DoesNotExists
	}
	if data.expiresAt.After(time.Now()) {
		delete(s.store, userID)
		return Expired
	}

	return nil
}
