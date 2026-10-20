package project_init

import (
	"fmt"
	"log"
	"os"

	"github.com/FatAndreasbot/go_project/user_service/infra/config"
)

func SetupApplictaion(options *LaunchOptions) error {

	err := config.LoadConfig(options.EnvPath)
	if err != nil {
		fmt.Println(os.Getwd())
		log.Default().Println(err)
	}

	userRepo, groupRepo, permissionRepo, session, err := SetupOutgoingAdapters()
	if err != nil {
		return err
	}

	jwtHandler, permissionHandler := setupBaseHandlers(
		userRepo,
		groupRepo,
	)

	groupRequests, userRequests, loginRequests := setupHandlers(
		userRepo, session, groupRepo, permissionRepo,
		jwtHandler, permissionHandler,
	)

	return setupIncomingAdapters(
		groupRequests, userRequests, loginRequests,
		session, jwtHandler,
	)
}
