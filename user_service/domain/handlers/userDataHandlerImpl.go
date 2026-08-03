package handlers

import (
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

// GetAndCheckUserByUsername(string) (*models.User, error)
// StoreNewUser(string, string, uuid.UUID) (*models.User, error)
// GetGroupList(pagesize, pagenumber int) ([]*models.Group, error)

// implementing incoming.IncomingRequestHandler
func (h *UserDataHandler) GetAndCheckUserByUsername(username, password string) (*models.User, error) {
	user, err := h.userPersistanceAdapter.GetUserByUsername(username)
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
func (h *UserDataHandler) StoreNewUser(username, password string, groupID uuid.UUID) (*models.User, error) {
	group, err := h.groupPersistanceAdapter.GetGroupByID(groupID)

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

	userUUID, err := h.userPersistanceAdapter.CreateUser(&user)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}
	user.ID = userUUID

	return &user, nil
}

func (h *UserDataHandler) GetGroupList(pagesize, pagenumber int) ([]*models.Group, error) {
	start := (pagenumber - 1) * pagesize
	end := start + pagesize + 1

	return h.groupPersistanceAdapter.GetGroupList(start, end)
}
