package application

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type JWTHandlerPort interface {
	EncodeAccessJWT(ctx context.Context, sessionID uuid.UUID) (string, time.Time, error)
	EncodeRefreshJWT(ctx context.Context, sessionID uuid.UUID) (string, time.Time, error)
	DecodeJWT(ctx context.Context, signedToken string) (uuid.UUID, error)
}
