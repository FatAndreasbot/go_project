package commands

import "fmt"

func Dispatch(args ...string) {
	if args[0] == "createsuperuser" {
		createsuperuser()
	} else {
		fmt.Println("unknown command")
	}
}
