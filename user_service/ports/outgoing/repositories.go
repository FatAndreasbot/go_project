package outgoing

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByEmail(context.Context, string) (*models.User, error)
	GetUserByID(context.Context, uuid.UUID) (*models.User, error)
	GetUsers(context.Context, models.Cursor, int) ([]*models.User, error)
	CreateUser(ctx context.Context, userData *models.UserData, pwordHash string) (*models.User, error)
	UpdateUser(context.Context, *models.UserUpdateRequestData) error
	CheckIdempotencyKey(context.Context, uuid.UUID) error
	DeleteUser(context.Context, uuid.UUID) error
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, groupName string, permissionIDs uuid.UUIDs, idempotencyKey uuid.UUID) (*models.Group, error)
	GetGroupByName(context.Context, string) (*models.MainGroupData, error)
	GetGroupByID(context.Context, uuid.UUID) (*models.MainGroupData, error)
	GetGroups(context.Context, models.Cursor, int) ([]*models.MainGroupData, error)
	UpdateGroup(context.Context, *models.MainGroupData) error
	DeleteGroup(context.Context, uuid.UUID) error
	GetGroupPermissions(context.Context, uuid.UUID) ([]*models.Permission, error)
	AddPermissions(ctx context.Context, groupID uuid.UUID, permissionIDs uuid.UUIDs) error
	RemovePermissions(ctx context.Context, groupID uuid.UUID, permissionIDs uuid.UUIDs) error

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
