package sessionstorage

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type MapSessionStorage struct {
	store map[uuid.UUID]any
	mu    sync.Mutex
}

func NewMapSessionStorate() *MapSessionStorage {
	return &MapSessionStorage{
		store: make(map[uuid.UUID]any),
		mu:    sync.Mutex{},
	}
}

func (s *MapSessionStorage) CreateSession(ctx context.Context, userID uuid.UUID) error
func (s *MapSessionStorage) DeleteSession(ctx context.Context, userID uuid.UUID) error
func (s *MapSessionStorage) ExtendSession(ctx context.Context, userID uuid.UUID) error
