package project_init

import (
	"fmt"
	"log"
	"os"

	"github.com/FatAndreasbot/go_project/user_service/infra/config"
)

func SetupApplictaion(options *LaunchOptions) error {

	config, err := config.LoadConfig(options.EnvPath)
	if err != nil {
		fmt.Println(os.Getwd())
		log.Default().Println(err)
	}

	userRepo, groupRepo, permissionRepo, session, err := SetupOutgoingAdapters(config)
	if err != nil {
		return err
	}

	jwtHandler, permissionHandler, passwordEncrypter := SetupBaseHandlers(
		config,
		userRepo,
		groupRepo,
	)

	groupRequests, userRequests, loginRequests := SetupHandlers(
		userRepo, session, groupRepo, permissionRepo,
		jwtHandler, permissionHandler, passwordEncrypter,
	)

	return setupIncomingAdapters(
		config,
		groupRequests, userRequests, loginRequests,
		session, permissionHandler, jwtHandler,
	)
}
