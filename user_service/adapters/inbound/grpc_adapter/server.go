package grpcadapter

import (
	"context"
	"errors"
	"log"
	proto "proto/user_service/v1"

	dominerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	handler        incoming.IncomingRequestHandler
	sessionHandler incoming.SessionStateHandler
}

func NewServer(handler incoming.IncomingRequestHandler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s Server) LogIn(ctx context.Context, req *proto.LogInRequest) (*proto.LogInResponse, error) {
	email, password := req.GetEmail(), req.GetPassword()

	user, err := s.handler.CheckUser(ctx, email, password)
	if err != nil {
		if errors.Is(err, dominerrors.WrongPasswdOrNoUserFound) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	access, err := EncodeAccessJWT(user)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	refresh, err := EncodeRefreshJWT(user)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	err = s.sessionHandler.AddSession(ctx, user.ID)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &proto.LogInResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		UserId:       user.ID.String(),
	}, nil
}

func (s Server) LogOut(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	rawUserID := ctx.Value(userdataKey)
	if rawUserID == nil {
		log.Default().Println("attempted to log out without an access token")
		return &emptypb.Empty{}, nil
	}
	userID, ok := rawUserID.(uuid.UUID)
	if !ok {
		log.Default().Println("the userID was not a UUID")
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	s.sessionHandler.RevokeSession(ctx, userID)
	return &emptypb.Empty{}, nil
}

func (s Server) Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	// refreshToken := req.GetRefreshToken()
	// data, err := DecodeJWT(refreshToken)

	return nil, status.Error(codes.Unimplemented, "method CreateUser not implemented")
}

func (Server) CreateUser(context.Context, *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateUser not implemented")
}
func (Server) GetUser(context.Context, *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetUser not implemented")
}
func (Server) GetUsers(context.Context, *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetUsers not implemented")
}
func (Server) GetPrivateUserData(context.Context, *proto.GetUserRequest) (*proto.GetPrivateUserDataResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetPrivateUserData not implemented")
}
func (Server) UpdateUser(context.Context, *proto.UpdateUserDataRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateUser not implemented")
}
func (Server) DeleteUser(context.Context, *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteUser not implemented")
}
func (Server) CreateGroup(context.Context, *proto.CreateGroupRequest) (*proto.CreateGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateGroup not implemented")
}
func (Server) GetGroup(context.Context, *proto.GetGroupRequest) (*proto.GetGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetGroup not implemented")
}
func (Server) GetGroups(context.Context, *proto.GetGroupsRequest) (*proto.GetGroupsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetGroups not implemented")
}
func (Server) UpdateGroup(context.Context, *proto.UpdateGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateGroup not implemented")
}
func (Server) DeleteGroup(context.Context, *proto.DeleteGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteGroup not implemented")
}
func (Server) GetPermissions(context.Context, *emptypb.Empty) (*proto.GetPermissionsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetPermissions not implemented")
}
