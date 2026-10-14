package incoming

import (
	"context"

	"github.com/google/uuid"
)

type SessionStateHandler interface {
	AddSession(ctx context.Context, userID uuid.UUID) error
	RevokeSession(ctx context.Context, userID uuid.UUID) error
	UpdateSession(ctx context.Context, userID string)
}
