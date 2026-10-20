package incoming

import (
	"context"
	"errors"
	"time"

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
	CreateUser(ctx context.Context, userdata *models.UserData, requesterID uuid.UUID) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID, requesterID uuid.UUID) (*models.User, error)
	GetUsers(ctx context.Context, cursor time.Time, count int, requesterID uuid.UUID) (users []*models.User, next time.Time, err error)
	UpdateUser(ctx context.Context, data *models.UserUpdateRequestData, requesterID uuid.UUID) error
	DeleteUser(ctx context.Context, userToDeleteID uuid.UUID, requesterID uuid.UUID) error
}

type GroupRequestHandler interface {
	CreateGroup(ctx context.Context, groupName string, permissions []*models.Permission, idempotencyKey, requesterID uuid.UUID) (*models.Group, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID, requesterID uuid.UUID) (*models.Group, error)
	GetGroups(ctx context.Context, cursor time.Time, count int, requesterID uuid.UUID) (groups []*models.MainGroupData, next time.Time, err error)
	UpdateGroup(ctx context.Context, groupID uuid.UUID, newGroupName string, reqesterID uuid.UUID) error
	AddPermission(ctx context.Context, groupID uuid.UUID, permission *models.Permission, requesterID uuid.UUID) error
	RemovePermission(ctx context.Context, groupID uuid.UUID, permission *models.Permission, requesterID uuid.UUID) error
	DeleteGroup(ctx context.Context, groupID uuid.UUID, requesterID uuid.UUID) error
	GetAllPermissions(ctx context.Context) ([]*models.Permission, error)
}
