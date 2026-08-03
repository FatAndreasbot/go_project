package grpcadapter

import (
	"context"
	"errors"
	"log"
	v1 "proto/common/v1"
	proto "proto/user_service"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/domain/models/dominaerrors"
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

func convertPermissions(domainPermissions *[]*models.Permission) []*v1.UserPermission {
	rpcPermissions := make([]*v1.UserPermission, 0, len(*domainPermissions))

	for _, permission := range *domainPermissions {
		rpcPermissions = append(rpcPermissions, &v1.UserPermission{
			PermissionId: &v1.UUID{Id: permission.ID.String()},
			Name:         permission.Name,
		})
	}
	return rpcPermissions
}

func (s *Server) LogIn(ctx context.Context, req *proto.LogInRequest) (*proto.LogInResponse, error) {
	username, password := req.GetUsername(), req.GetPassword()

	user, err := s.handler.GetAndCheckUserByUsername(username, password)
	if err != nil {
		if errors.Is(err, dominaerrors.UserNotFoundError) {
			return &proto.LogInResponse{Success: false}, status.Error(codes.NotFound, "user not found")
		} else if errors.Is(err, dominaerrors.UserWrongPassword) {
			return &proto.LogInResponse{Success: false}, status.Error(codes.PermissionDenied, "wrong password")
		} else {
			return &proto.LogInResponse{Success: false}, status.Error(codes.Internal, "could not fetch userdata")
		}
	}

	jwt, err := EncodeJWT(user.ID)
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

	user, err := s.handler.StoreNewUser(
		req.GetUsername(),
		req.GetPassword(),
		groupUUID,
	)
	permissions := convertPermissions(&user.Group.Permissions)

	return &proto.CreateUserResponse{
		User: &v1.User{
			UserId: &v1.UUID{Id: user.ID.String()},
			Name:   user.Name,
			Group: &v1.UserGroup{
				GroupId:         &v1.UUID{Id: user.Group.ID.String()},
				GroupName:       user.Group.Name,
				UserPermissions: permissions,
			},
		},
	}, nil
}

func (s *Server) GetGroups(ctx context.Context, req *proto.GetGroupsRequest) (*proto.GetGroupsResponse, error) {
	groups, err := s.handler.GetGroupList(int(req.GetPagesize()), int(req.GetPagenumber()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	convertedGroups := make([]*v1.UserGroup, 0, int(req.GetPagesize()))

	for _, group := range groups {
		convertedGroups = append(convertedGroups, &v1.UserGroup{
			GroupId:         &v1.UUID{Id: group.ID.String()},
			GroupName:       group.Name,
			UserPermissions: convertPermissions(&group.Permissions),
		})
	}

	return &proto.GetGroupsResponse{
		Groups: convertedGroups,
	}, nil
}
