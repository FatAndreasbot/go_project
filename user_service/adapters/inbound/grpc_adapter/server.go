package grpcadapter

import (
	proto "proto/user_service/v1"

	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	loginHandler incoming.LoginHandler
	userHandler  incoming.UserRequestHandler
	groupHandler incoming.GroupRequestHandler
}

func NewServer(handler incoming.LoginHandler) *Server {
	return &Server{
		loginHandler: handler,
	}
}
