package outgoing

import (
	"context"

	"github.com/google/uuid"
)

type SessionHandler interface {
	CreateSession(ctx context.Context) (uuid.UUID, error)
	DeleteSession(ctx context.Context, sessionKey uuid.UUID) error
	ExtendSession(ctx context.Context, sessionKey uuid.UUID) error
	CheckSession(ctx context.Context, sessionKey uuid.UUID) error

	SetData(ctx context.Context, sessionKey uuid.UUID, key string, data any) error
	GetData(ctx context.Context, sessionKey uuid.UUID, key string) (any, error)
}
