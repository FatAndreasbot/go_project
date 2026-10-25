package application

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type PermissionHandlerPort interface {
	// True if has permission
	HasPermission(ctx context.Context, userID uuid.UUID, permission models.PermissionName) bool
}
