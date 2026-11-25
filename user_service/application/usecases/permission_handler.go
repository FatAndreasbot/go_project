package usecases

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PermissionHandler struct {
	userRepository  outgoing.UserRepository
	groupRepository outgoing.GroupRepository
}

func NewPermissionHander(
	userRepository outgoing.UserRepository,
	groupRepository outgoing.GroupRepository,
) *PermissionHandler {
	return &PermissionHandler{
		userRepository:  userRepository,
		groupRepository: groupRepository,
	}
}

func (h *PermissionHandler) HasPermission(ctx context.Context, userID uuid.UUID, permission models.PermissionName) bool {
	if permission == models.Any {
		return true
	}

	userData, err := h.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		zap.L().Error(err.Error())
		return false
	}

	groupPermissions, err := h.groupRepository.GetGroupPermissions(ctx, userData.GroupID)
	if err != nil {
		zap.L().Error(err.Error())
		return false
	}

	for _, groupPermission := range groupPermissions {
		if groupPermission.Name == permission {
			return true
		}
	}

	return false
}
