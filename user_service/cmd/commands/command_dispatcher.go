package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/FatAndreasbot/go_project/user_service/infra/config"
)

func Dispatch(args ...string) {
	err := config.LoadConfig("../.env")
	if err != nil {
		fmt.Println(os.Getwd())
		log.Default().Println(err)
	}

	if args[0] == "createsuperuser" {
		createsuperuser()
	} else {
		fmt.Println("unknown command")
	}
}
