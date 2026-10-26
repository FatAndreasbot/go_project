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

type UserRequestsHandlerImpl struct {
	userRepository    outgoing.UserRepository
	permissionHandler application.PermissionHandlerPort
}

func NewUserRequestHandler(
	userRepository outgoing.UserRepository,
	permissionHandler application.PermissionHandlerPort,
) *UserRequestsHandlerImpl {
	return &UserRequestsHandlerImpl{
		userRepository:    userRepository,
		permissionHandler: permissionHandler,
	}
}

func (h *UserRequestsHandlerImpl) CreateUser(ctx context.Context, userdata *models.UserData) (*models.User, error) {
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

func (h *UserRequestsHandlerImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := h.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, DBError
	}
	return user, nil
}

func (h *UserRequestsHandlerImpl) GetUsers(ctx context.Context, cursor models.Cursor, count int) ([]*models.User, models.Cursor, error) {
	users, err := h.userRepository.GetUsers(
		ctx,
		cursor,
		count,
	)

	if err != nil {
		log.Default().Println(err.Error())
		return nil, models.Cursor{}, DBError
	}

	if len(users) == 0 {
		return nil, models.Cursor{}, domainerrors.EmptyList
	}
	var next time.Time
	var nextID uuid.UUID
	next = users[0].CreatedAt
	nextID = users[0].ID

	// if for some accursed reason the list is not sorted by CreatedAt
	for _, user := range users {
		if user.CreatedAt.After(next) {
			next = user.CreatedAt
			nextID = user.ID
		}
	}

	return users, models.Cursor{
		Timestamp: next,
		ID:        nextID,
	}, nil
}

func (h *UserRequestsHandlerImpl) UpdateUser(ctx context.Context, data *models.UserUpdateRequestData) error {
	if err := h.userRepository.UpdateUser(ctx, data); err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	return nil
}

func (h *UserRequestsHandlerImpl) DeleteUser(ctx context.Context, userToDeleteID uuid.UUID) error {
	if err := h.userRepository.DeleteUser(ctx, userToDeleteID); err != nil {
		log.Default().Println(err.Error())
		return DBError
	}

	return nil
}
