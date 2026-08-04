package incoming

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type IncomingRequestHandler interface {
	GetAndCheckUserByUsername(ctx context.Context, username, password string) (*models.User, error)
	StoreNewUser(ctx context.Context, username, password string, groupID uuid.UUID) (*models.User, error)
	GetGroupList(ctx context.Context, pagesize, pagenumber int) ([]*models.Group, error)
}
