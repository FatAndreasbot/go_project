package project_init

import "flag"

type LaunchOptions struct {
	EnvPath    string
	AppCommand string
	IsDebug    bool
}

var options *LaunchOptions

func ParseLaunchFlags() *LaunchOptions {
	envFile := flag.String("env-file", "../.env", "Path to .env file")
	command := flag.String("command", "serve", "Action this software should do")
	isDebug := flag.Bool("is-debug", false, "To set up debug mode")

	flag.Parse()

	options = &LaunchOptions{
		EnvPath:    *envFile,
		AppCommand: *command,
		IsDebug:    *isDebug,
	}

	return options
}
