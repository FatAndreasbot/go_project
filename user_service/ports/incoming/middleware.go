package incoming

import (
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
)

type Middleware interface {
	GetAuthFunction([]string) models.Interceptor
	GetPermissionFunction(map[string]models.PermissionName, []string) models.Interceptor
}
