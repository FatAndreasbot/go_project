package main

import (
	"fmt"

	"github.com/FatAndreasbot/go_project/user_service/cmd/project_init"
)

func main() {
	err := project_init.SetupApplication()
	if err != nil {
		fmt.Println(err)
	}
}
