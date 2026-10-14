package grpcadapter

import (
	"context"
	"errors"
	"log"
	v1 "proto/common/v1"
	proto "proto/user_service/v1"

	dominaerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	requestHandler incoming.IncomingRequestHandler
	// sessionHandler incoming.SessionStateHandler
	// jwtHandler     incoming.JWTHandler
}

func NewServer(handler incoming.IncomingRequestHandler) *Server {
	return &Server{
		requestHandler: handler,
	}
}

func (s Server) LogIn(ctx context.Context, req *proto.LogInRequest) (*proto.LogInResponse, error) {
	email, password := req.GetEmail(), req.GetPassword()

	repsponseData, err := s.requestHandler.Login(ctx, email, password)
	if err != nil {
		if errors.Is(err, dominaerrors.InvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "wromg email or password")
		}
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &proto.LogInResponse{
		AccessToken:  repsponseData.AccessToken,
		RefreshToken: repsponseData.RefreshToken,
		UserId:       repsponseData.UserID.String(),
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

	s.requestHandler.Logout(ctx, userID)
	return &emptypb.Empty{}, nil
}

func (s Server) Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	refreshToken := req.GetRefreshToken()

	tokens, err := s.requestHandler.Refresh(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, dominaerrors.InvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "invalid refresh token")
		}
		if errors.Is(err, dominaerrors.ExpiredCredentials) {
			return nil, status.Error(codes.PermissionDenied, "expired refresh token")
		}
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &proto.RefreshResponse{
		NewAccessToken:  tokens.AccessToken,
		NewRefreshToken: tokens.RefreshToken,
	}, nil
}

func (s Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	groupIDString := req.GetGroupId()

	groupID, err := uuid.Parse(groupIDString)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	userdata := models.UserData{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		GroupID:  groupID,
	}

	user, err := s.requestHandler.CreateUser(ctx, &userdata, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &proto.CreateUserResponse{
		User: &v1.UserData{
			Id:      user.ID.String(),
			Email:   user.Email,
			GroupId: user.GroupID.String(),
		},
	}, nil
}

func (s Server) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	userIDString := req.GetUserId()
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	user, err := s.requestHandler.GetUserByID(ctx, userID, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &proto.GetUserResponse{
		User: &v1.UserData{
			Email:   user.Email,
			GroupId: user.GroupID.String(),
			Id:      user.ID.String(),
		},
	}, nil
}

func (s Server) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	cursor := req.GetCursor().AsTime()
	count := int(req.Limit)
	users, next, err := s.requestHandler.GetUsers(ctx, cursor, count, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	resp := proto.GetUsersResponse{
		Cursor: timestamppb.New(next),
		Users:  make([]*v1.UserData, 0, count),
	}

	for _, user := range users {
		resp.Users = append(resp.Users, &v1.UserData{
			Id:      user.ID.String(),
			Email:   user.Email,
			GroupId: user.GroupID.String(),
		})
	}

	return &resp, nil
}

func (s Server) UpdateUser(ctx context.Context, req *proto.UpdateUserDataRequest) (*emptypb.Empty, error) {
	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	userToChangeID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "expected UUID")
	}

	newGroupID, err := uuid.Parse(req.GetUserData().GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "expected UUID")
	}

	err = s.requestHandler.UpdateUser(ctx, &models.UserUpdateRequestData{
		UserToChangeID: userToChangeID,
		NewPassword:    req.GetUserData().GetPassword(),
		NewEmail:       req.GetUserData().GetEmail(),
		NewGroupID:     newGroupID,
	}, requesterID)

	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &emptypb.Empty{}, nil
}

func (s Server) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	// requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	// if !ok {
	// 	return nil, status.Error(codes.Internal, "unexpected error")
	// }

	return nil, status.Error(codes.Unimplemented, "method DeleteUser not implemented")
}

func (s Server) CreateGroup(ctx context.Context, req *proto.CreateGroupRequest) (*proto.CreateGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateGroup not implemented")
}

func (s Server) GetGroup(ctx context.Context, req *proto.GetGroupRequest) (*proto.GetGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetGroup not implemented")
}

func (s Server) GetGroups(ctx context.Context, req *proto.GetGroupsRequest) (*proto.GetGroupsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetGroups not implemented")
}

func (s Server) UpdateGroup(ctx context.Context, req *proto.UpdateGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateGroup not implemented")
}

func (s Server) DeleteGroup(ctx context.Context, req *proto.DeleteGroupRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteGroup not implemented")
}

func (s Server) GetPermissions(ctx context.Context, req *emptypb.Empty) (*proto.GetPermissionsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetPermissions not implemented")
}
