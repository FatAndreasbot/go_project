package outgoing

import (
	"context"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByEmail(context.Context, string) (*models.User, error)
	GetUserByID(context.Context, uuid.UUID) (*models.User, error)
	GetUsers(context.Context, time.Time, int) ([]*models.User, error)
	CreateUser(context.Context, *models.UserData) (*models.User, error)
	UpdateUser(context.Context, *models.UserUpdateRequestData) error
	CheckIdempotencyKey(context.Context, uuid.UUID) error
	DeleteUser(context.Context, uuid.UUID) error
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, groupName string, permissions []*models.Permission, idempotencyKey uuid.UUID) (*models.Group, error)
	GetGroupByID(context.Context, uuid.UUID) (*models.MainGroupData, error)
	GetGroups(context.Context, time.Time, int) ([]*models.MainGroupData, error)
	UpdateGroup(context.Context, *models.MainGroupData) error
	DeleteGroup(context.Context, uuid.UUID) error
	GetGroupPermissions(context.Context, uuid.UUID) ([]*models.Permission, error)
	AddPermission(ctx context.Context, groupID, permissionID uuid.UUID) error
	RemovePermission(ctx context.Context, groupID, permissionID uuid.UUID) error

	CheckIdempotencyKey(context.Context, uuid.UUID) error
}

type PermissionRepository interface {
	GetPermissionByID(context.Context, uuid.UUID) (*models.Permission, error)
	GetPermissions(context.Context) ([]*models.Permission, error)

	// no need for these methods
	// CreatePermission(context.Context, string) (*models.Permission, error)
	// UpdatePermission(context.Context, *models.Permission) error
	// DeletePermission(context.Context, uuid.UUID) error
}
