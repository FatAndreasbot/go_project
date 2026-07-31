package controllers

import (
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing/persistance"
	"github.com/google/uuid"
)

var groupController GroupController
var groupControllerWasInitialized bool = false

type GroupController struct {
	groupPersistanceAdapter persistance.GroupPersistancePort
}

func InitGroupController(adapter persistance.GroupPersistancePort) {
	groupController = GroupController{
		groupPersistanceAdapter: adapter,
	}
	groupControllerWasInitialized = true
}

func GetGroupController() *GroupController {
	if !groupControllerWasInitialized {
		GracefulShutdown()
		return nil
	}

	return &groupController
}

func (gc *GroupController) GetGroupByUUID(id uuid.UUID) (*models.Group, error) {
	return gc.groupPersistanceAdapter.GetGroupByID(id)
}

func (gc *GroupController) GetAllGroups() ([]*models.Group, error) {
	// TODO set actual limits
	return gc.groupPersistanceAdapter.GetGroupList(0, 1000)
}
