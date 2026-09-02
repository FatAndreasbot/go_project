package grpcadapter

import (
	"context"
	"errors"
	"log"
	v1 "proto/common/v1"
	proto "proto/user_service/v1"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/domain/models/dominaerrors"
	"github.com/FatAndreasbot/go_project/user_service/ports/incoming"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
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

func convertPermissions(domainPermissions *[]*models.Permission) []*v1.UserPermission {
	rpcPermissions := make([]*v1.UserPermission, 0, len(*domainPermissions))

	for _, permission := range *domainPermissions {
		rpcPermissions = append(rpcPermissions, &v1.UserPermission{
			PermissionId: permission.ID.String(),
			Name:         permission.Name,
		})
	}
	return rpcPermissions
}

func (s *Server) LogIn(ctx context.Context, req *proto.LogInRequest) (*proto.LogInResponse, error) {
	username, password := req.GetUsername(), req.GetPassword()

	user, err := s.handler.GetAndCheckUserByUsername(ctx, username, password)
	if err != nil {
		if errors.Is(err, dominaerrors.WrongPasswdOrNoUserFound) {
			return nil, status.Error(codes.NotFound, "wrong password user not found")
		} else {
			return nil, status.Error(codes.Internal, "could not fetch userdata")
		}
	}

	accessJWT, err := EncodeAccessJWT(user)
	if err != nil {
		log.Default().Println(err)
		return nil, status.Error(codes.Internal, "could not generate token")
	}
	refreshJWT, err := EncodeRefreshJWT(user)
	if err != nil {
		log.Default().Println(err)
		return nil, status.Error(codes.Internal, "could not generate token")
	}

	return &proto.LogInResponse{
		AccessToken: accessJWT,
		RefreshToken: refreshJWT,
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	groupUUID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		log.Default().Println(err)
		return nil, status.Error(codes.InvalidArgument, "could not parse groupID")
	}

	user, err := s.handler.StoreNewUser(
		ctx,
		req.GetUsername(),
		req.GetPassword(),
		groupUUID,
	)

	return &proto.CreateUserResponse{
		User: &v1.User{
			Id:  user.ID.String(),
			Name:    user.Name,
			GroupId: user.Group.ID.String(),
		},
	}, nil
}

func (s *Server) GetGroups(ctx context.Context, req *proto.GetGroupsRequest) (*proto.GetGroupsResponse, error) {
	groups, err := s.handler.GetGroupList(ctx, int(req.GetCursor()), int(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	convertedGroups := make([]*v1.UserGroup, 0, int(req.GetLimit()))

	for _, group := range groups {
		convertedGroups = append(convertedGroups, &v1.UserGroup{
			Id:         group.ID.String(),
			Name:       group.Name,
			UserPermissions: convertPermissions(&group.Permissions),
		})
	}

	return &proto.GetGroupsResponse{
		Groups: convertedGroups,
	}, nil
}

func (s *Server) GetPermissions(ctx context.Context, req *proto.GetPermissionsRequest) (*proto.GetPermissionsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "could not parse user id")
	}
	user, err := s.handler.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissions := convertPermissions(&user.Group.Permissions)

	return &proto.GetPermissionsResponse{
		UserPermissions: permissions,
	}, nil
}
