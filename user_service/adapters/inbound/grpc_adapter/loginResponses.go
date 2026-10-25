package grpcadapter

import (
	"context"
	"errors"
	"log"
	proto "proto/user_service/v1"

	"github.com/google/uuid"

	dominaerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LogIn(ctx context.Context, req *proto.LogInRequest) (*proto.LogInResponse, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	email, password := req.GetEmail(), req.GetPassword()

	repsponseData, err := s.loginHandler.Login(ctx, email, password)
	if err != nil {
		if errors.Is(err, dominaerrors.InvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "wromg email or password")
		}
		return nil, status.Error(codes.Internal, "unexpected error")
	}

	return &proto.LogInResponse{
		AccessToken:      repsponseData.AccessToken,
		RefreshToken:     repsponseData.RefreshToken,
		UserId:           repsponseData.UserID.String(),
		AccessExpiresAt:  timestamppb.New(repsponseData.AccesExpiresAt),
		RefreshExpiresAt: timestamppb.New(repsponseData.RefreshExpiresAt),
	}, nil
}

func (s *Server) LogOut(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

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

	s.loginHandler.Logout(ctx, userID)
	return &emptypb.Empty{}, nil
}

func (s *Server) Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	if !s.listening.Load() {
		return nil, status.Error(codes.Unavailable, "server is shutting down")
	}
	s.requestsInProgressCounter.Add(1)
	defer s.requestsInProgressCounter.Add(-1)

	refreshToken := req.GetRefreshToken()

	tokens, err := s.loginHandler.Refresh(ctx, refreshToken)
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
