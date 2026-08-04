package project_init

import (
	"github.com/FatAndreasbot/go_project/user_service/domain/handlers"
	"github.com/FatAndreasbot/go_project/user_service/infra/config"
)

func SetupApplication() error {
	err := config.InitEnvVars()
	if err != nil {
		return err
	}

	userAdp, groupAdp, err := setupOutgoingAdapters()
	if err != nil {
		return err
	}

	userDataHandler := handlers.NewUserDataHandler(userAdp, groupAdp)

	return setupIncomingAdapters(userDataHandler)
}
