package handlers

import (
	"context"
	"errors"
	"log"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"github.com/google/uuid"
)

var (
	DBError      = errors.New("DB query error")
	JWTError     = errors.New("token encoding/decoding error")
	SessionError = errors.New("session error")
)

// implements ports/incoming/request.go LoginHandler
type LoginRequestsHandler struct {
	userRepository outgoing.UserRepository
	jwtHandler     outgoing.JWTHandler
	sessionHandler outgoing.SessionStateHandler
}

func NewRequestHandler() *LoginRequestsHandler

func (h *LoginRequestsHandler) Login(ctx context.Context, email, password string) (*models.LoginResponseData, error) {
	user, err := h.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		user.CheckPassword(password)
		log.Default().Println(err.Error())
		if errors.Is(err, domainerrors.NotFound) {
			return nil, domainerrors.InvalidCredentials
		}
		return nil, DBError
	}

	err = user.CheckPassword(password)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, domainerrors.InvalidCredentials
	}

	accessToken, err := h.jwtHandler.EncodeAccessJWT(ctx, user)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, JWTError
	}
	refreshToken, err := h.jwtHandler.EncodeRefreshJWT(ctx, user)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, JWTError
	}

	err = h.sessionHandler.CreateSession(ctx, user.ID)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, SessionError
	}

	return &models.LoginResponseData{
		UserID: user.ID,
		Tokens: models.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (h *LoginRequestsHandler) Logout(ctx context.Context, userID uuid.UUID) error {
	return h.sessionHandler.DeleteSession(ctx, userID)
}

func (h *LoginRequestsHandler) Refresh(ctx context.Context, RefreshToken string) (*models.Tokens, error) {
	userID, err := h.jwtHandler.DecodeJWT(ctx, RefreshToken)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, err
	}
	err = h.sessionHandler.ExtendSession(ctx, userID)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, err
	}

	user, err := h.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domainerrors.NotFound) {
			return nil, domainerrors.InvalidCredentials
		}
		return nil, DBError
	}

	NewAccessToken, err := h.jwtHandler.EncodeAccessJWT(ctx, user)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, JWTError
	}
	newRefreshToken, err := h.jwtHandler.EncodeRefreshJWT(ctx, user)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, JWTError
	}

	return &models.Tokens{
		AccessToken:  NewAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
