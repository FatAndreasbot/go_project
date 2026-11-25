package main

import (
	"github.com/FatAndreasbot/go_project/user_service/cmd/commands"
	"github.com/FatAndreasbot/go_project/user_service/cmd/project_init"
)

func main() {
	options := project_init.ParseLaunchFlags()
	project_init.SetupLogger(options)

	commands.Dispatch(options)
}
