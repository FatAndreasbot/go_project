package project_init

import (
	"fmt"
	"net"
	proto "proto/user_service/v1"

	grpcadapter "github.com/FatAndreasbot/go_project/user_service/adapters/inbound/grpc_adapter"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"google.golang.org/grpc"
)

func setupIncomingAdapters(
	config *models.ConfigOptions,

	groupHandler incoming.GroupRequestHandler,
	userHandler incoming.UserRequestHandler,
	loginHandler incoming.LoginRequestHandler,

	sessionHandler outgoing.SessionHandler,

	jwtHandler application.JWTHandlerPort,
) error {
	port := config.GetPort()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	interceptor := grpcadapter.NewAuthInterceptor(
		jwtHandler,
		sessionHandler,
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()),
	)
	serverAdapter := grpcadapter.NewServer(
		groupHandler, userHandler, loginHandler,
	)

	proto.RegisterUserServiceServer(grpcServer, serverAdapter)

	setupTerminationCallback(func() {
		grpcServer.GracefulStop()
		// serverAdapter.Shutdown()
		listener.Close()
	})

	fmt.Printf("starting server at localhost:%d\n", port)
	return grpcServer.Serve(listener)
}
