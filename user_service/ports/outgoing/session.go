package outgoing

import (
	"context"

	"github.com/google/uuid"
)

type SessionHandler interface {
	CreateSession(ctx context.Context, userID uuid.UUID) error
	DeleteSession(ctx context.Context, userID uuid.UUID) error
	ExtendSession(ctx context.Context, userID uuid.UUID) error
	CheckSession(ctx context.Context, userID uuid.UUID) error
}
