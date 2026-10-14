package grpcadapter

import (
	"context"
	"errors"
	"log"
	v1 "proto/common/v1"
	proto "proto/user_service/v1"

	dominaerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

	user, err := s.userHandler.CreateUser(ctx, &userdata, requesterID)
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

	user, err := s.userHandler.GetUserByID(ctx, userID, requesterID)
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
	users, next, err := s.userHandler.GetUsers(ctx, cursor, count, requesterID)
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

	err = s.userHandler.UpdateUser(ctx, &models.UserUpdateRequestData{
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
	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	userToDeleteID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "expected UUID")
	}

	err = s.userHandler.DeleteUser(ctx, userToDeleteID, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &emptypb.Empty{}, nil
}
