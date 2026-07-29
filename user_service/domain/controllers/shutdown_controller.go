package controllers

var shutdownAdapterActions []func() = make([]func(), 0)

func GracefulShutdown() {
	for _, shutdown := range shutdownAdapterActions {
		shutdown()
	}
}

func RegisterShutdownAction(action func()) {
	shutdownAdapterActions = append(shutdownAdapterActions, action)
}
