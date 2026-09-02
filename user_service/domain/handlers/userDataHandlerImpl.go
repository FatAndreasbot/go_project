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

func NewUserDataHandler(userAdp persistance.UserPersistancePort, groupAdp persistance.GroupPersistancePort) *UserDataHandler {
	return &UserDataHandler{
		userPersistanceAdapter:  userAdp,
		groupPersistanceAdapter: groupAdp,
	}
}

// implementing incoming.IncomingRequestHandler
func (h *UserDataHandler) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := h.userPersistanceAdapter.GetUserByUsername(ctx, username)
	if err != nil {
		log.Default().Println(errors.Join(err, dominaerrors.WrongPasswdOrNoUserFound))
		return nil, dominaerrors.WrongPasswdOrNoUserFound
	}

	return user, nil

}

// implementing incoming.IncomingRequestHandler
func (h *UserDataHandler) StoreNewUser(ctx context.Context, username, password string, groupID, userID uuid.UUID) (*models.User, error) {
	group, err := h.groupPersistanceAdapter.GetGroupByID(ctx, groupID)

	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("group by id %q was not found", groupID.String())
	}

	user := models.User{
		ID:    userID,
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

func (h *UserDataHandler) GetGroupList(ctx context.Context, pagesize int, lastGroupName string) ([]*models.Group, error) {
	return h.groupPersistanceAdapter.GetGroupList(ctx, pagesize, lastGroupName)
}

func (h *UserDataHandler) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return h.userPersistanceAdapter.GetUserByID(ctx, userID)
}
