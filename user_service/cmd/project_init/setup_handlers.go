package project_init

import (
	"github.com/FatAndreasbot/go_project/user_service/application/handlers"
	"github.com/FatAndreasbot/go_project/user_service/application/handlers/requests"
	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
)

func setupBaseHandlers(
	userRepository outgoing.UserRepository,
	groupRepository outgoing.GroupRepository,
) (
	jwtHandler application.JWTHandler,
	permissionHandler application.PermissionHandler,
) {
	jwt := handlers.NewJWTHandler()
	permissions := handlers.NewPermissionHander(
		userRepository,
		groupRepository,
	)

	return jwt, permissions
}

func setupHandlers(
	userRepository outgoing.UserRepository,
	sessionHandler outgoing.SessionHandler,
	groupRepository outgoing.GroupRepository,
	permissionRepository outgoing.PermissionRepository,

	jwtHandler application.JWTHandler,
	permissionHandler application.PermissionHandler,
) (
	incoming.GroupRequestHandler,
	incoming.UserRequestHandler,
	incoming.LoginRequestHandler,
) {
	loginHandler := requests.NewLoginRequestHandler(
		userRepository, sessionHandler, jwtHandler,
	)
	groupHandler := requests.NewGroupRequestHandler(
		groupRepository, permissionRepository, permissionHandler,
	)
	userHandler := requests.NewUserRequestHandler(
		userRepository, permissionHandler,
	)

	return groupHandler, userHandler, loginHandler
}
