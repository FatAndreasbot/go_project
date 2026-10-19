package project_init

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func SetupTerminationCallback(callback func()) {
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh

	fmt.Println("shutting down server...")
	callback()
	fmt.Println("server shut down")

	signal.Stop(sigCh)
}
