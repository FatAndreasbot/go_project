package grpcadapter

import (
	proto "proto/user_service/v1"

	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	handler incoming.IncomingRequestHandler
}

func NewServer(handler incoming.IncomingRequestHandler) *Server {
	return &Server{
		handler: handler,
	}
}
