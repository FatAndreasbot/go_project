package project_init

import (
	"fmt"
	"log"
	"os"

	"github.com/FatAndreasbot/go_project/user_service/infra/config"
	"go.uber.org/zap"
)

func SetupLogger(options *LaunchOptions) {
	config, err := config.LoadConfig(options.EnvPath)
	if err != nil {
		fmt.Println(os.Getwd())
		log.Default().Println(err)
	}

	var logger *zap.Logger

	if config.GetDebug() {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}

	zap.ReplaceGlobals(logger)
}
