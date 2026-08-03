package incoming

import (
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type IncomingRequestHandler interface {
	GetAndCheckUserByUsername(username, password string) (*models.User, error)
	StoreNewUser(username, password string, groupID uuid.UUID) (*models.User, error)
	GetGroupList(pagesize, pagenumber int) ([]*models.Group, error)
}
