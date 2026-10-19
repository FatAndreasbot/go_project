package requests

import (
	"context"
	"log"
	"time"

	"slices"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"github.com/google/uuid"
)

type GroupRequestHandler struct {
	groupRepository      outgoing.GroupRepository
	permissionRepository outgoing.PermissionRepository

	permissionHandler application.PermissionHandler
}

func NewGroupRequestHandler(
	groupRepository outgoing.GroupRepository,
	permissionRepository outgoing.PermissionRepository,

	permissionHandler application.PermissionHandler,
) *GroupRequestHandler {
	return &GroupRequestHandler{
		groupRepository:      groupRepository,
		permissionRepository: permissionRepository,

		permissionHandler: permissionHandler,
	}
}

func (h *GroupRequestHandler) CreateGroup(ctx context.Context, groupName string, permissions []*models.Permission, idempotencyKey, requesterID uuid.UUID) (*models.Group, error) {
	if !h.permissionHandler.HasPermission(requesterID, models.GroupCreate) {
		return nil, domainerrors.NotEnoughPermission
	}
	err := h.groupRepository.CheckIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, domainerrors.IdempotencyFail
	}

	newGroup, err := h.groupRepository.CreateGroup(
		ctx,
		groupName,
		permissions,
		idempotencyKey,
	)

	if err != nil {
		log.Default().Println(err.Error())
		return nil, DBError
	}

	return newGroup, nil
}

func (h *GroupRequestHandler) GetGroupByID(ctx context.Context, groupID uuid.UUID, requesterID uuid.UUID) (*models.Group, error) {
	if !h.permissionHandler.HasPermission(requesterID, models.GroupRead) {
		return nil, domainerrors.NotEnoughPermission
	}
	groupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, DBError
	}

	permissions, err := h.groupRepository.GetGroupPermissions(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, DBError
	}

	group := models.Group{
		MainGroupData: *groupData,
		Permissions:   permissions,
	}

	return &group, nil
}

func (h *GroupRequestHandler) GetGroups(ctx context.Context, cursor time.Time, count int, requesterID uuid.UUID) ([]*models.MainGroupData, time.Time, error) {
	if !h.permissionHandler.HasPermission(requesterID, models.GroupRead) {
		return nil, time.Time{}, domainerrors.NotEnoughPermission
	}
	groups, err := h.groupRepository.GetGroups(ctx, cursor, count)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, time.Time{}, DBError
	}

	// next := groups[len(groups)-1].CreatedAt
	var next time.Time
	next = groups[0].CreatedAt

	// if for some accursed reason the list is not sorted by CreatedAt
	for _, group := range groups {
		if group.CreatedAt.After(next) {
			next = group.CreatedAt
		}
	}

	return groups, next, nil
}

func (h *GroupRequestHandler) UpdateGroup(ctx context.Context, groupID uuid.UUID, newGroupName string, requesterID uuid.UUID) error {
	if !h.permissionHandler.HasPermission(requesterID, models.GroupUpdate) {
		return domainerrors.NotEnoughPermission
	}

	groupToEdit, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return DBError
	}
	groupToEdit.Name = newGroupName

	if err := h.groupRepository.UpdateGroup(ctx, groupToEdit); err != nil {
		log.Default().Println(err.Error())
		return err
	}

	return nil
}

func (h *GroupRequestHandler) AddPermission(ctx context.Context, groupID uuid.UUID, permission *models.Permission, requesterID uuid.UUID) error {
	if !h.permissionHandler.HasPermission(requesterID, models.GroupUpdate) {
		return domainerrors.NotEnoughPermission
	}

	groupPermissions, err := h.groupRepository.GetGroupPermissions(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	if slices.ContainsFunc(groupPermissions, func(p *models.Permission) bool {
		return *p == *permission
	}) {
		return GroupHasPermission
	}

	if err := h.groupRepository.AddPermission(ctx, groupID, permission.ID); err != nil {
		log.Default().Println(err.Error())
		return err
	}

	return nil
}

func (h *GroupRequestHandler) RemovePermission(ctx context.Context, groupID uuid.UUID, permission *models.Permission, requesterID uuid.UUID) error {
	if !h.permissionHandler.HasPermission(requesterID, models.GroupUpdate) {
		return domainerrors.NotEnoughPermission
	}

	groupPermissions, err := h.groupRepository.GetGroupPermissions(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	cmpFunc := func(p *models.Permission) bool {
		return *p == *permission
	}
	if !slices.ContainsFunc(groupPermissions, cmpFunc) {
		return GroupLacksPermission
	}

	if err := h.groupRepository.RemovePermission(ctx, groupID, permission.ID); err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	return nil
}

func (h *GroupRequestHandler) DeleteGroup(ctx context.Context, groupID uuid.UUID, requesterID uuid.UUID) error {
	if !h.permissionHandler.HasPermission(requesterID, models.GroupDelete) {
		return domainerrors.NotEnoughPermission
	}

	if err := h.groupRepository.DeleteGroup(ctx, groupID); err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	return domainerrors.NotImplemented
}

func (h *GroupRequestHandler) GetAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	permissions, err := h.permissionRepository.GetPermissions(ctx)

	if err != nil {
		log.Default().Println(err.Error())
		return nil, DBError
	}

	return permissions, nil
}
