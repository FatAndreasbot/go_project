package handlers

import (
	"context"
	"log"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"github.com/google/uuid"
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

func (h *PermissionHandler) HasPermission(userID uuid.UUID, permission models.PermissionName) bool {
	ctx := context.Background()

	userData, err := h.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		log.Default().Println(err)
		return false
	}

	groupPermissions, err := h.groupRepository.GetGroupPermissions(ctx, userData.GroupID)
	if err != nil {
		log.Default().Println(err)
		return false
	}

	for _, groupPermission := range groupPermissions {
		if groupPermission.Name == permission {
			return true
		}
	}

	return false
}
