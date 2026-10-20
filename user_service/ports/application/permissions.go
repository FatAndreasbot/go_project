package application

import (
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type PermissionHandlerPort interface {
	// True if has permission
	HasPermission(userID uuid.UUID, permission models.PermissionName) bool
}
