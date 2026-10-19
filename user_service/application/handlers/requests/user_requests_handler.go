package requests

import (
	"context"
	"log"
	"time"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"github.com/google/uuid"
)

type UserRequestsHandler struct {
	userRepository    outgoing.UserRepository
	permissionHandler application.PermissionHandler
}

func (h *UserRequestsHandler) CreateUser(ctx context.Context, userdata *models.UserData, requesterID uuid.UUID) (*models.User, error) {
	if !h.permissionHandler.HasPermission(requesterID, models.UserCreate) {
		return nil, domainerrors.NotEnoughPermission
	}

	err := h.userRepository.CheckIdempotencyKey(ctx, userdata.IdempotencyKey)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, domainerrors.IdempotencyFail
	}

	user, err := h.userRepository.CreateUser(ctx, userdata)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, DBError
	}

	return user, nil
}

func (h *UserRequestsHandler) GetUserByID(ctx context.Context, userID uuid.UUID, requesterID uuid.UUID) (*models.User, error) {
	if userID != requesterID && !h.permissionHandler.HasPermission(requesterID, models.UserRead) {
		return nil, domainerrors.NotEnoughPermission
	}
	user, err := h.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, DBError
	}
	return user, nil
}

func (h *UserRequestsHandler) GetUsers(ctx context.Context, cursor time.Time, count int, requesterID uuid.UUID) ([]*models.User, time.Time, error) {
	if !h.permissionHandler.HasPermission(requesterID, models.UserCreate) {
		return nil, time.Time{}, domainerrors.NotEnoughPermission
	}

	users, err := h.userRepository.GetUsers(
		ctx,
		cursor,
		count,
	)

	if err != nil {
		log.Default().Println(err.Error())
		return nil, time.Time{}, DBError
	}

	var next time.Time
	next = users[0].CreatedAt

	// if for some accursed reason the list is not sorted by CreatedAt
	for _, group := range users {
		if group.CreatedAt.After(next) {
			next = group.CreatedAt
		}
	}

	return users, next, nil
}

func (h *UserRequestsHandler) UpdateUser(ctx context.Context, data *models.UserUpdateRequestData, requesterID uuid.UUID) error {
	if !h.permissionHandler.HasPermission(requesterID, models.UserUpdate) {
		return domainerrors.NotEnoughPermission
	}

	if err := h.userRepository.UpdateUser(ctx, data); err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	return nil
}

func (h *UserRequestsHandler) DeleteUser(ctx context.Context, userToDeleteID uuid.UUID, requesterID uuid.UUID) error {
	if !h.permissionHandler.HasPermission(requesterID, models.UserDelete) {
		return domainerrors.NotEnoughPermission
	}
	if err := h.userRepository.DeleteUser(ctx, userToDeleteID); err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	return nil
}
