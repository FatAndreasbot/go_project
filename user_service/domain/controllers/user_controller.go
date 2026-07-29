package controllers

import (
	"fmt"
	"log"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing/persistance"
	"github.com/google/uuid"
)

var userController UserController
var userControllerWasInitialized bool = false

type UserController struct {
	userPersistanceAdapter persistance.UserPersistancePort
}

func InitUserController(adapter persistance.UserPersistancePort) {
	userController = UserController{
		userPersistanceAdapter: adapter,
	}
	userControllerWasInitialized = true
}

func GetUserController() *UserController {
	if !userControllerWasInitialized {
		GracefulShutdown()
		return nil
	}

	return &userController
}

func (uc *UserController) CreateUser(name, password string, groupID uuid.UUID) (*models.User, error) {
	gc := GetGroupController()
	group, err := gc.GetGroupByUUID(groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("group by id %q was not found", groupID.String())
	}

	user := models.User{
		Name:  name,
		Group: group,
	}
	err = user.SetPassword(password)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}

	userUUID, err := uc.userPersistanceAdapter.CreateUser(&user)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}
	user.ID = userUUID

	return &user, nil
}
