package commands

import (
	"fmt"
	"os"

	"github.com/FatAndreasbot/go_project/user_service/cmd/project_init"
)

func Dispatch(command *project_init.LaunchOptions) {
	if command.IsDebug {
		fmt.Printf("PID: %d\n ", os.Getpid())
	}

	switch command.AppCommand {
	case "createsuperuser":
		createsuperuser(command)
	case "serve":
		err := project_init.SetupApplictaion(command)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("unknown command")
	}
}
