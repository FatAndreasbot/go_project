package main

import (
	"fmt"
	"os"

	"github.com/FatAndreasbot/go_project/user_service/cmd/commands"
	"github.com/FatAndreasbot/go_project/user_service/cmd/project_init"
)

func main() {
	if len(os.Args) == 1 {
		err := project_init.SetupApplictaion()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		commands.Dispatch(os.Args[1:]...)
	}
}
