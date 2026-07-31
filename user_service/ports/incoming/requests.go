package incoming

import (
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type IncomingRequestHandler interface {
	GetUserByUsername(string) (*models.User, error)
	GetJWT(*models.User) (string, error)
	StoreNewUser(string, string, uuid.UUID) (*models.User, error)
	GetGroupList(pagesize, pagenumber int) ([]*models.Group, error)
}
