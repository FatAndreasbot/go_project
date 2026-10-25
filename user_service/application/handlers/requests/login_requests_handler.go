package requests

import (
	"context"
	"errors"
	"log"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"

	"github.com/google/uuid"
)

// implements ports/incoming/request.go LoginHandler
type LoginRequestsHandlerImpl struct {
	userRepository outgoing.UserRepository
	sessionHandler outgoing.SessionHandler

	jwtHandler application.JWTHandlerPort
}

func NewLoginRequestHandler(
	userRepository outgoing.UserRepository,
	sessionHandler outgoing.SessionHandler,

	jwtHandler application.JWTHandlerPort,
) *LoginRequestsHandlerImpl {
	return &LoginRequestsHandlerImpl{
		userRepository: userRepository,
		sessionHandler: sessionHandler,

		jwtHandler: jwtHandler,
	}
}

func (h *LoginRequestsHandlerImpl) Login(ctx context.Context, email, password string) (*models.LoginResponseData, error) {
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

	accessToken, accessExpiresAt, err := h.jwtHandler.EncodeAccessJWT(ctx, user)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, JWTError
	}
	refreshToken, refreshExpiresAt, err := h.jwtHandler.EncodeRefreshJWT(ctx, user)
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
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			AccesExpiresAt:   accessExpiresAt,
			RefreshExpiresAt: refreshExpiresAt,
		},
	}, nil
}

func (h *LoginRequestsHandlerImpl) Logout(ctx context.Context, userID uuid.UUID) error {
	return h.sessionHandler.DeleteSession(ctx, userID)
}

func (h *LoginRequestsHandlerImpl) Refresh(ctx context.Context, RefreshToken string) (*models.Tokens, error) {
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
		log.Default().Println(err.Error())
		if errors.Is(err, domainerrors.NotFound) {
			return nil, domainerrors.InvalidCredentials
		}
		return nil, DBError
	}

	NewAccessToken, accessExpiresAt, err := h.jwtHandler.EncodeAccessJWT(ctx, user)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, JWTError
	}
	newRefreshToken, refreshExpiresAt, err := h.jwtHandler.EncodeRefreshJWT(ctx, user)
	if err != nil {
		log.Default().Println(err.Error())
		return nil, JWTError
	}

	return &models.Tokens{
		AccessToken:      NewAccessToken,
		RefreshToken:     newRefreshToken,
		AccesExpiresAt:   accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}
