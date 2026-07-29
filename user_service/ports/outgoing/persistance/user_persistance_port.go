package persistance

import (
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type UserPersistancePort interface {
	GetUserByUsername(string) (*models.User, error)
	GetUserByID(uuid.UUID) (*models.User, error)

	GetUserList(int, int) ([]*models.User, error)

	UpdateUser(uuid.UUID, *models.User) error
	CreateUser(*models.User) (uuid.UUID, error)
	DeleteUser(uuid.UUID) error
}

type GroupPersistancePort interface {
	GetGroupByID(uuid.UUID) (*models.Group, error)

	GetGroupList(int, int) ([]*models.Group, error)

	UpdateGroup(uuid.UUID, *models.Group) error
	CreateGroup(*models.Group) (uuid.UUID, error)
	DeleteGroup(uuid.UUID) error
}
