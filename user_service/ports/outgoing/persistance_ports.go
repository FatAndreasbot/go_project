package outgoing

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByEmail(context.Context, string) (*models.User, error)
	GetUserByID(context.Context, uuid.UUID) (*models.User, error)
}

type GroupPersistancePort interface {
	GetGroupByID(context.Context, uuid.UUID) (*models.Group, error)

	GetGroupList(ctx context.Context, limit int, name string) ([]*models.Group, error)

	UpdateGroup(context.Context, uuid.UUID, *models.Group) error
	CreateGroup(context.Context, *models.Group) (uuid.UUID, error)
	DeleteGroup(context.Context, uuid.UUID) error
}
