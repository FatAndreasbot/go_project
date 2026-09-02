package incoming

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type IncomingRequestHandler interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	StoreNewUser(ctx context.Context, username, password string, groupID uuid.UUID) (*models.User, error)
	GetGroupList(ctx context.Context, limit int, lastGroupName string) ([]*models.Group, error)
	GetUserByID(ctx context.Context, userid uuid.UUID) (*models.User, error)
}
