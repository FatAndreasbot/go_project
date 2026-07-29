package incoming

import "github.com/FatAndreasbot/go_project/user_service/domain/models"

type IncomingRequestHandler interface {
	GetUserByUsername(string) (*models.User, error)
	GetJWT(*models.User) (string, error)
	StoreNewUser(*models.User) error
}
