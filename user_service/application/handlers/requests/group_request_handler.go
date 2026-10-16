package requests

import (
	"context"
	"time"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

type GroupRequestHandler struct{}

func (h GroupRequestHandler) CreateGroup(ctx context.Context, groupName string, permissions []*models.Permission, idempotencyKey, requesterID uuid.UUID) (*models.Group, error) {
	return nil, domainerrors.NotImplemented
}

func (h GroupRequestHandler) GetGroupByID(ctx context.Context, groupID uuid.UUID, requesterID uuid.UUID) (*models.Group, error) {
	return nil, domainerrors.NotImplemented
}

func (h GroupRequestHandler) GetGroups(ctx context.Context, cursor time.Time, count int, requesterID uuid.UUID) (groups []*models.Group, next time.Time, err error) {
	return nil, time.Time{}, domainerrors.NotImplemented
}

func (h GroupRequestHandler) UpdateGroup(ctx context.Context, groupID uuid.UUID, newGroupName string, reqesterID uuid.UUID) error {
	return domainerrors.NotImplemented
}

func (h GroupRequestHandler) AddPermission(ctx context.Context, groupID uuid.UUID, permission *models.Permission, requesterID uuid.UUID) error {
	return domainerrors.NotImplemented
}

func (h GroupRequestHandler) RemovePermission(ctx context.Context, groupID uuid.UUID, permission *models.Permission, requesterID uuid.UUID) error {
	return domainerrors.NotImplemented
}

func (h GroupRequestHandler) DeleteGroup(ctx context.Context, groupID uuid.UUID, requesterID uuid.UUID) error {
	return domainerrors.NotImplemented
}

func (h GroupRequestHandler) GetAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	return nil, domainerrors.NotImplemented
}
