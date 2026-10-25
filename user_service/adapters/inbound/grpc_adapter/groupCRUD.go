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

func domainPermossionsToProtoPermissions(domianPermissions []*models.Permission) []*v1.UserPermission {
	permissions := make([]*v1.UserPermission, 0, len(domianPermissions))

	for _, domianPermission := range domianPermissions {
		permissions = append(permissions, &v1.UserPermission{
			PermissionId: domianPermission.ID.String(),
			Name:         string(domianPermission.Name),
		})
	}

	return permissions
}

func protoPermissionToDomainPermissions(protoPermissions []*v1.UserPermission) ([]*models.Permission, error) {
	permissions := make([]*models.Permission, 0, len(protoPermissions))

	for _, protoPermission := range protoPermissions {
		permissionID, err := uuid.Parse(protoPermission.GetPermissionId())
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &models.Permission{
			ID:   permissionID,
			Name: models.PermissionName(protoPermission.GetName()),
		})
	}

	return permissions, nil
}

func (s *Server) CreateGroup(ctx context.Context, req *proto.CreateGroupRequest) (*proto.CreateGroupResponse, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	protoPermissions := req.GetPermissions()
	permissions, err := protoPermissionToDomainPermissions(protoPermissions)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	idempotencyKey, err := uuid.Parse(req.GetIdempotencyKey())
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	group, err := s.groupHandler.CreateGroup(ctx, req.GetName(), permissions, idempotencyKey, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	response := proto.CreateGroupResponse{
		Group: &v1.UserGroup{
			Id:              group.ID.String(),
			Name:            group.Name,
			UserPermissions: protoPermissions,
		},
	}

	return &response, status.Error(codes.Unimplemented, "method CreateGroup not implemented")
}

func (s *Server) GetGroup(ctx context.Context, req *proto.GetGroupRequest) (*proto.GetGroupResponse, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	groupID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	group, err := s.groupHandler.GetGroupByID(ctx, groupID, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	response := proto.GetGroupResponse{
		Group: &v1.UserGroup{
			Id:              group.ID.String(),
			Name:            group.Name,
			UserPermissions: domainPermossionsToProtoPermissions(group.Permissions),
		},
	}

	return &response, nil
}

func (s *Server) GetGroups(ctx context.Context, req *proto.GetGroupsRequest) (*proto.GetGroupsResponse, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	cursor := req.GetCursor().AsTime()
	count := int(req.Limit)

	groups, next, err := s.groupHandler.GetGroups(ctx, cursor, count, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	resp := proto.GetGroupsResponse{
		Cursor: timestamppb.New(next),
		Groups: make([]*v1.UserGroupMainData, 0, count),
	}

	for _, group := range groups {
		resp.Groups = append(resp.Groups, &v1.UserGroupMainData{
			Id:   group.ID.String(),
			Name: group.Name,
		})
	}

	return &resp, nil
}

func (s *Server) UpdateGroup(ctx context.Context, req *proto.UpdateGroupRequest) (*emptypb.Empty, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	groupID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	err = s.groupHandler.UpdateGroup(ctx, groupID, req.GetNewName(), requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) AddPermissionToGroup(ctx context.Context, req *proto.AddPermissionToGroupRequest) (*emptypb.Empty, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	groupID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	permissionID, err := uuid.Parse(req.GetPermissionId())
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	err = s.groupHandler.AddPermission(ctx, groupID, permissionID, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) RemovePermissionFromGroup(ctx context.Context, req *proto.RemovePermissionFromGroupRequest) (*emptypb.Empty, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	groupID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	permissionID, err := uuid.Parse(req.GetPermission().GetPermissionId())
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	permission := models.Permission{
		ID:   permissionID,
		Name: models.PermissionName(req.GetPermission().GetName()),
	}

	err = s.groupHandler.RemovePermission(ctx, groupID, &permission, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteGroup(ctx context.Context, req *proto.DeleteGroupRequest) (*emptypb.Empty, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	requesterID, ok := ctx.Value(userdataKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "unexpected error")
	}
	groupID, err := uuid.Parse(req.GetGroupId())
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.InvalidArgument, "expected UUID")
	}

	err = s.groupHandler.DeleteGroup(ctx, groupID, requesterID)
	if err != nil {
		if errors.Is(err, dominaerrors.NotEnoughPermission) {
			return nil, status.Error(codes.PermissionDenied, "user lacks permissions")
		}
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetAllPermissions(ctx context.Context, req *emptypb.Empty) (*proto.GetPermissionsResponse, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	permissions, err := s.groupHandler.GetAllPermissions(ctx)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	permissionsProto := domainPermossionsToProtoPermissions(permissions)

	return &proto.GetPermissionsResponse{
		Permissions: permissionsProto,
	}, nil
}
