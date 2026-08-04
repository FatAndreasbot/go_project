package commands

import "errors"

var CommnadNotFoundError error = errors.New("command not found")
var NotImplementedError error = errors.New("command not implemented")

func DispatchCommand(command string) error {
	switch command {
	case "seed_db":
		return seedDB()
	default:
		return CommnadNotFoundError
	}
}
