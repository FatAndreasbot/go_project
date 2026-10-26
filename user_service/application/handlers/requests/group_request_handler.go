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

func (h *GroupRequestHandlerImpl) CreateGroup(ctx context.Context, groupName string, permissions []*models.Permission, idempotencyKey uuid.UUID) (*models.Group, error) {
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

func (h *GroupRequestHandlerImpl) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error) {
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

func (h *GroupRequestHandlerImpl) GetGroups(ctx context.Context, cursor models.Cursor, count int) ([]*models.MainGroupData, models.Cursor, error) {
	groups, err := h.groupRepository.GetGroups(ctx, cursor, count)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, models.Cursor{}, DBError
	}

	if len(groups) == 0 {
		return nil, models.Cursor{}, domainerrors.EmptyList
	}
	// next := groups[len(groups)-1].CreatedAt
	var next time.Time
	var nextID uuid.UUID
	next = groups[0].CreatedAt
	nextID = groups[0].ID

	// if for some accursed reason the list is not sorted by CreatedAt
	for _, group := range groups {
		if group.CreatedAt.After(next) {
			next = group.CreatedAt
			nextID = group.ID
		}
	}

	return groups, models.Cursor{
		Timestamp: next,
		ID:        nextID,
	}, nil
}

func (h *GroupRequestHandlerImpl) UpdateGroup(ctx context.Context, groupID uuid.UUID, newGroupName string) error {
	oldGroupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return DBError
	}
	if oldGroupData.Name == "admin" {
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

func (h *GroupRequestHandlerImpl) AddPermission(ctx context.Context, groupID uuid.UUID, permissionID uuid.UUID) error {
	oldGroupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return DBError
	}
	if oldGroupData.Name == "admin" {
		return domainerrors.NotEnoughPermission
	}

	groupPermissions, err := h.groupRepository.GetGroupPermissions(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	if slices.ContainsFunc(groupPermissions, func(p *models.Permission) bool {
		return p.ID == permissionID
	}) {
		return GroupHasPermission
	}

	if err := h.groupRepository.AddPermission(ctx, groupID, permissionID); err != nil {
		log.Default().Println(err.Error())
		return err
	}

	return nil
}

func (h *GroupRequestHandlerImpl) RemovePermission(ctx context.Context, groupID uuid.UUID, permissionID uuid.UUID) error {
	oldGroupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return DBError
	}
	if oldGroupData.Name == "admin" {
		return domainerrors.NotEnoughPermission
	}

	groupPermissions, err := h.groupRepository.GetGroupPermissions(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	cmpFunc := func(p *models.Permission) bool {
		return p.ID == permissionID
	}
	if !slices.ContainsFunc(groupPermissions, cmpFunc) {
		return GroupLacksPermission
	}

	if err := h.groupRepository.RemovePermission(ctx, groupID, permissionID); err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	return nil
}

func (h *GroupRequestHandlerImpl) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	oldGroupData, err := h.groupRepository.GetGroupByID(ctx, groupID)
	if err != nil {
		log.Default().Println(err.Error())
		return DBError
	}
	if oldGroupData.Name == "admin" {
		return domainerrors.NotEnoughPermission
	}

	if err := h.groupRepository.DeleteGroup(ctx, groupID); err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	return domainerrors.NotImplemented
}

func (h *GroupRequestHandlerImpl) GetAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	permissions, err := h.permissionRepository.GetPermissions(ctx)

	if err != nil {
		log.Default().Println(err.Error())
		return nil, DBError
	}

	return permissions, nil
}
