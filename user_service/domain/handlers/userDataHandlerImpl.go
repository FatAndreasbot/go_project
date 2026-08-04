package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/domain/models/dominaerrors"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing/persistance"
	"github.com/google/uuid"
)

type UserDataHandler struct {
	userPersistanceAdapter  persistance.UserPersistancePort
	groupPersistanceAdapter persistance.GroupPersistancePort
}

// implementing incoming.IncomingRequestHandler
func (h *UserDataHandler) GetAndCheckUserByUsername(ctx context.Context, username, password string) (*models.User, error) {
	user, err := h.userPersistanceAdapter.GetUserByUsername(ctx, username)
	if err != nil {
		log.Default().Println(errors.Join(err, dominaerrors.UserNotFoundError))
		return nil, dominaerrors.UserNotFoundError
	}

	err = user.CheckPassword(password)
	if err != nil {
		log.Default().Println(err)
		return nil, dominaerrors.UserWrongPassword
	}

	return user, nil

}

// implementing incoming.IncomingRequestHandler
func (h *UserDataHandler) StoreNewUser(ctx context.Context, username, password string, groupID uuid.UUID) (*models.User, error) {
	group, err := h.groupPersistanceAdapter.GetGroupByID(ctx, groupID)

	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("group by id %q was not found", groupID.String())
	}

	user := models.User{
		Name:  username,
		Group: group,
	}
	err = user.SetPassword(password)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}

	userUUID, err := h.userPersistanceAdapter.CreateUser(ctx, &user)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}
	user.ID = userUUID

	return &user, nil
}

func (h *UserDataHandler) GetGroupList(ctx context.Context, pagesize, pagenumber int) ([]*models.Group, error) {
	offset := (pagenumber - 1) * pagesize

	return h.groupPersistanceAdapter.GetGroupList(ctx, pagesize, offset)
}
