package grpcadapter

import (
	proto "proto/user_service/v1"
	"sync/atomic"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	loginHandler incoming.LoginRequestHandler
	userHandler  incoming.UserRequestHandler
	groupHandler incoming.GroupRequestHandler

	requestsInProgressCounter atomic.Int32
	listening                 atomic.Bool
}

func NewServer(
	groupHandler incoming.GroupRequestHandler,
	userHandler incoming.UserRequestHandler,
	loginHandler incoming.LoginRequestHandler,
) *Server {
	s := Server{
		loginHandler: loginHandler,
		userHandler:  userHandler,
		groupHandler: groupHandler,
	}

	s.listening.Store(true)
	s.requestsInProgressCounter.Store(0)

	return &s
}

func (s *Server) Shutdown() {
	s.listening.Store(false)

	for s.requestsInProgressCounter.Load() > 0 {
		time.Sleep(time.Millisecond * 200)
	}
}
