package project_init

import (
	"github.com/FatAndreasbot/go_project/user_service/application/handlers"
	"github.com/FatAndreasbot/go_project/user_service/application/handlers/requests"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
)

func setupBaseHandlers(
	config *models.ConfigOptions,
	userRepository outgoing.UserRepository,
	groupRepository outgoing.GroupRepository,
) (
	jwtHandler application.JWTHandlerPort,
	permissionHandler application.PermissionHandlerPort,
	passwordEncryptionHandler application.PasswordEncrypter,
) {
	jwt := handlers.NewJWTHandler(config)
	permissions := handlers.NewPermissionHander(
		userRepository,
		groupRepository,
	)

	passwordEncrypter := handlers.BcryptPasswordEncrypter{}

	return jwt, permissions, passwordEncrypter
}

func setupHandlers(
	userRepository outgoing.UserRepository,
	sessionHandler outgoing.SessionHandler,
	groupRepository outgoing.GroupRepository,
	permissionRepository outgoing.PermissionRepository,

	jwtHandler application.JWTHandlerPort,
	permissionHandler application.PermissionHandlerPort,
	passwordEncrypter application.PasswordEncrypter,
) (
	incoming.GroupRequestHandler,
	incoming.UserRequestHandler,
	incoming.LoginRequestHandler,
) {
	loginHandler := requests.NewLoginRequestHandler(
		userRepository, sessionHandler, jwtHandler, passwordEncrypter,
	)
	groupHandler := requests.NewGroupRequestHandler(
		groupRepository, permissionRepository, permissionHandler,
	)
	userHandler := requests.NewUserRequestHandler(
		userRepository, permissionHandler,
	)

	return groupHandler, userHandler, loginHandler
}
