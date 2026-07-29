package grpcadapter

import (
	"context"
	"log"
	proto "proto/user_service"

	"github.com/FatAndreasbot/go_project/user_service/domain/controllers"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

// LogIn(context.Context, *LogInRequest) (*LogInResponse, error)
// CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
// GetGroups(context.Context, *GetGroupsRequest) (*GetGroupsResponse, error)

func (s *Server) LogIn(ctx context.Context, req *proto.LogInRequest) (*proto.LogInResponse, error) {
	username, password := req.GetUsername(), req.GetPassword()

	user, err := s.handler.GetUserByUsername(username)
	if err != nil {
		log.Default().Println(err)
		return &proto.LogInResponse{Success: false}, status.Error(codes.NotFound, "user not found")
	}

	err = user.CheckPassword(password)
	if err != nil {
		log.Default().Println(err)
		return &proto.LogInResponse{Success: false}, status.Error(codes.PermissionDenied, "wrong password")
	}

	jwt, err := s.handler.GetJWT(user)
	if err != nil {
		log.Default().Println(err)
		return &proto.LogInResponse{Success: false}, status.Error(codes.Internal, "could not generate token")
	}

	metadata.AppendToOutgoingContext(
		ctx,
		"authorization",
		"Bearer "+jwt,
	)

	return &proto.LogInResponse{Success: true}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	groupUUID, err := uuid.Parse(req.GetGroupId().GetId())
	if err != nil {
		log.Default().Println(err)
		return nil, status.Error(codes.InvalidArgument, "could not parse groupID")
	}

	user, err := controllers.GetUserController().CreateUser(
		req.GetUsername(),
		req.GetPassword(),
		groupUUID,
	)
	if 

	return nil, status.Error(codes.Unimplemented, "not implemented")
}
