package project_init

import (
	"fmt"
	"net"
	proto "proto/user_service"

	grpcadapter "github.com/FatAndreasbot/go_project/user_service/adapters/inbound/grpc_adapter"
	"github.com/FatAndreasbot/go_project/user_service/infra/config"
	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
)

func setupIncomingAdapters(handler incoming.IncomingRequestHandler) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetNetworkPort()))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor(grpcadapter.Authenticate)),
	)

	proto.RegisterUserServiceServer(grpcServer, grpcadapter.NewServer(handler))
	return grpcServer.Serve(listener)
}
