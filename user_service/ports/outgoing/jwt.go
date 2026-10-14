package outgoing

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type JWTHandler interface {
	EncodeAccessJWT(ctx context.Context, user *models.User) (string, error)
	EncodeRefreshJWT(ctx context.Context, user *models.User) (string, error)
	DecodeJWT(ctx context.Context, signedToken string) (uuid.UUID, error)
}
