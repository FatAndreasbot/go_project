package project_init

import (
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
)

func SetupOutgoingAdapters() (
	outgoing.UserRepository,
	outgoing.GroupRepository,
	outgoing.PermissionRepository,

	outgoing.SessionHandler,
)

// {
// 	upa := postgresrepository.NewUserPersistanceAdapter()
// 	gpa := postgresrepository.NewUserPersistanceAdapter()
// 	ppa
// }
