package persistance

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type UserPersistancePort interface {
	GetUserByUsername(context.Context, string) (*models.User, error)
	GetUserByID(context.Context, uuid.UUID) (*models.User, error)

	GetUserList(ctx context.Context, limit, offset int) ([]*models.User, error)

	UpdateUser(context.Context, uuid.UUID, *models.User) error
	CreateUser(context.Context, *models.User) (uuid.UUID, error)
	DeleteUser(context.Context, uuid.UUID) error
}

type GroupPersistancePort interface {
	GetGroupByID(context.Context, uuid.UUID) (*models.Group, error)

	GetGroupList(ctx context.Context, limit, offset int) ([]*models.Group, error)

	UpdateGroup(context.Context, uuid.UUID, *models.Group) error
	CreateGroup(context.Context, *models.Group) (uuid.UUID, error)
	DeleteGroup(context.Context, uuid.UUID) error
}
