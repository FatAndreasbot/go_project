package application

import (
	"context"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type JWTHandlerPort interface {
	EncodeAccessJWT(ctx context.Context, user *models.User) (string, time.Time, error)
	EncodeRefreshJWT(ctx context.Context, user *models.User) (string, time.Time, error)
	DecodeJWT(ctx context.Context, signedToken string) (uuid.UUID, error)
}
