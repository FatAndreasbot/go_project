package requests

import (
	"context"
	"errors"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

// implements ports/incoming/request.go LoginHandler
type LoginRequestsHandlerImpl struct {
	userRepository outgoing.UserRepository
	sessionHandler outgoing.SessionHandler

	jwtHandler        application.JWTUsecasePort
	passwordEncrypter application.PasswordEncrypter
}

func NewLoginRequestHandler(
	userRepository outgoing.UserRepository,
	sessionHandler outgoing.SessionHandler,

	jwtHandler application.JWTUsecasePort,
	passwordEncrypter application.PasswordEncrypter,
) *LoginRequestsHandlerImpl {
	return &LoginRequestsHandlerImpl{
		userRepository: userRepository,
		sessionHandler: sessionHandler,

		jwtHandler:        jwtHandler,
		passwordEncrypter: passwordEncrypter,
	}
}

const userIDSessionKey string = "USER_ID"

func (h *LoginRequestsHandlerImpl) Login(ctx context.Context, email, password string) (*models.LoginResponseData, error) {
	user, err := h.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		h.passwordEncrypter.FakeCheckPassword(password)

		zap.L().Error(err.Error())
		if errors.Is(err, domainerrors.NotFound) {
			return nil, domainerrors.InvalidCredentials
		}
		return nil, DBError
	}

	err = h.passwordEncrypter.CheckPassword(password, user)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, domainerrors.InvalidCredentials
	}

	sessionKey, err := h.sessionHandler.CreateSession(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, SessionError
	}
	h.sessionHandler.SetData(ctx, sessionKey, userIDSessionKey, user.ID)

	accessToken, accessExpiresAt, err := h.jwtHandler.EncodeAccessJWT(ctx, sessionKey)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, JWTError
	}
	refreshToken, refreshExpiresAt, err := h.jwtHandler.EncodeRefreshJWT(ctx, sessionKey)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, JWTError
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
	sessionID, err := h.jwtHandler.DecodeJWT(ctx, RefreshToken)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	err = h.sessionHandler.ExtendSession(ctx, sessionID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	NewAccessToken, accessExpiresAt, err := h.jwtHandler.EncodeAccessJWT(ctx, sessionID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, JWTError
	}
	newRefreshToken, refreshExpiresAt, err := h.jwtHandler.EncodeRefreshJWT(ctx, sessionID)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, JWTError
	}

	return &models.Tokens{
		AccessToken:      NewAccessToken,
		RefreshToken:     newRefreshToken,
		AccesExpiresAt:   accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}
