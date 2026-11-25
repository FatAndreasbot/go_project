package requests

import (
	"context"
	"errors"
	"time"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GroupRequestHandlerImpl struct {
	groupRepository      outgoing.GroupRepository
	permissionRepository outgoing.PermissionRepository

	permissionHandler application.PermissionHandlerPort
}

func NewGroupRequestHandler(
	groupRepository outgoing.GroupRepository,
	permissionRepository outgoing.PermissionRepository,

	permissionHandler application.PermissionHandlerPort,
) *GroupRequestHandlerImpl {
	return &GroupRequestHandlerImpl{
		groupRepository:      groupRepository,
		permissionRepository: permissionRepository,

		permissionHandler: permissionHandler,
	}
}

func (h *GroupRequestHandlerImpl) CreateGroup(ctx context.Context, groupName string, permissionIDs uuid.UUIDs, idempotencyKey uuid.UUID) (*models.Group, error) {
	err := h.groupRepository.CheckIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, domainerrors.IdempotencyFail
	}

	newGroup, err := h.groupRepository.CreateGroup(
		ctx,
		groupName,
		permissionIDs,
		idempotencyKey,
	)

	if err != nil {
		zap.L().Error(err.Error())
		return nil, DBError
	}

	return newGroup, nil
}

func (h *GroupRequestHandlerImpl) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error) {
	groupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, DBError
	}

	permissions, err := h.groupRepository.GetGroupPermissions(ctx, groupID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, DBError
	}

	group := models.Group{
		MainGroupData: *groupData,
		Permissions:   permissions,
	}

	return &group, nil
}

func (h *GroupRequestHandlerImpl) GetGroups(ctx context.Context, cursor models.Cursor, count int) ([]*models.MainGroupData, models.Cursor, error) {
	groups, err := h.groupRepository.GetGroups(ctx, cursor, count)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, models.Cursor{}, DBError
	}

	if len(groups) == 0 {
		return nil, models.Cursor{}, domainerrors.EmptyList
	}

	var next time.Time
	var nextID uuid.UUID
	// assuming that groups is sorted by CreatedAt
	next = groups[len(groups)-1].CreatedAt
	nextID = groups[len(groups)-1].ID

	return groups, models.Cursor{
		Timestamp: next,
		ID:        nextID,
	}, nil
}

func (h *GroupRequestHandlerImpl) UpdateGroup(ctx context.Context, groupID uuid.UUID, newGroupName string) error {
	oldGroupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		zap.L().Error(err.Error())
		return DBError
	}
	if oldGroupData.Name == "admin" {
		return domainerrors.NotEnoughPermission
	}

	groupToEdit, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		zap.L().Error(err.Error())
		return DBError
	}
	groupToEdit.Name = newGroupName

	if err := h.groupRepository.UpdateGroup(ctx, groupToEdit); err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

func (h *GroupRequestHandlerImpl) AddPermissions(ctx context.Context, groupID uuid.UUID, permissionIDList uuid.UUIDs) error {
	oldGroupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		zap.L().Error(err.Error())
		return DBError
	}
	if oldGroupData.Name == "admin" {
		return domainerrors.NotEnoughPermission
	}

	if err := h.groupRepository.AddPermissions(ctx, groupID, permissionIDList); err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

func (h *GroupRequestHandlerImpl) RemovePermissions(ctx context.Context, groupID uuid.UUID, permissionIDs uuid.UUIDs) error {
	oldGroupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		zap.L().Error(err.Error())
		return DBError
	}
	if oldGroupData.Name == "admin" {
		return domainerrors.NotEnoughPermission
	}

	if err := h.groupRepository.RemovePermissions(ctx, groupID, permissionIDs); err != nil {
		zap.L().Error(err.Error())
		return DBError
	}

	return nil
}

func (h *GroupRequestHandlerImpl) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	oldGroupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		zap.L().Error(err.Error())
		return DBError
	}
	if oldGroupData.Name == "admin" {
		return domainerrors.NotEnoughPermission
	}

	if err := h.groupRepository.DeleteGroup(ctx, groupID); err != nil {
		zap.L().Error(err.Error())
		if errors.Is(err, domainerrors.DeletingNonEmptyGroup) {
			return err
		}
		return DBError
	}

	return nil
}

func (h *GroupRequestHandlerImpl) GetAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	permissions, err := h.permissionRepository.GetPermissions(ctx)

	if err != nil {
		zap.L().Error(err.Error())
		return nil, DBError
	}

	return permissions, nil
}
