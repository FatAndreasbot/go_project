package incoming

import (
	"context"
	"errors"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

var (
	InternalError = errors.New("unexpected error")
)

type LoginRequestHandler interface {
	Login(ctx context.Context, email, password string) (*models.LoginResponseData, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	Refresh(ctx context.Context, RefreshToken string) (*models.Tokens, error)
}

type UserRequestHandler interface {
	CreateUser(ctx context.Context, userdata *models.UserData) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUsers(ctx context.Context, cursor models.Cursor, count int) (users []*models.User, next models.Cursor, err error)
	UpdateUser(ctx context.Context, data *models.UserUpdateRequestData) error
	DeleteUser(ctx context.Context, userToDeleteID uuid.UUID) error
}

type GroupRequestHandler interface {
	CreateGroup(ctx context.Context, groupName string, permissionIDs uuid.UUIDs, idempotencyKey uuid.UUID) (*models.Group, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
	GetGroups(ctx context.Context, cursor models.Cursor, count int) (groups []*models.MainGroupData, next models.Cursor, err error)
	UpdateGroup(ctx context.Context, groupID uuid.UUID, newGroupName string) error
	AddPermissions(ctx context.Context, groupID uuid.UUID, permissionIDList uuid.UUIDs) error
	RemovePermissions(ctx context.Context, groupID uuid.UUID, permission uuid.UUIDs) error
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error
	GetAllPermissions(ctx context.Context) ([]*models.Permission, error)
}
