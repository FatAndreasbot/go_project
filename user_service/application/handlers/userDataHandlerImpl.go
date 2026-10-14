package handlers

import (
	persistance "github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
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
