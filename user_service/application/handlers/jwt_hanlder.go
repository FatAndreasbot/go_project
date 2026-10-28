package handlers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTHandler struct {
	config *models.ConfigOptions
}

func NewJWTHandler(config *models.ConfigOptions) *JWTHandler {
	return &JWTHandler{
		config: config,
	}
}

func (h *JWTHandler) EncodeAccessJWT(ctx context.Context, user *models.User) (string, time.Time, error) {
	tokenIssuedAt := time.Now()

	expTime := h.config.GetAccessJWTExpiration()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(tokenIssuedAt.Add(expTime)),
		IssuedAt:  jwt.NewNumericDate(tokenIssuedAt),
	})

	secret := h.config.GetHS256Secret()

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, tokenIssuedAt.Add(expTime), nil
}

func (h *JWTHandler) EncodeRefreshJWT(ctx context.Context, user *models.User) (string, time.Time, error) {
	tokenIssuedAt := time.Now()

	expTime := h.config.GetRefreshJWTExpiration()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(tokenIssuedAt.Add(expTime)),
		IssuedAt:  jwt.NewNumericDate(tokenIssuedAt),
	})

	secret := h.config.GetHS256Secret()

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, tokenIssuedAt.Add(expTime), nil
}

func (h *JWTHandler) DecodeJWT(ctx context.Context, signedToken string) (uuid.UUID, error) {
	secret := h.config.GetHS256Secret()

	decoded, err := jwt.Parse(
		signedToken,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		log.Default().Println(err)
		return uuid.Nil, err
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)
	if !ok {
		log.Default().Printf("could not cast %v into MapClaims", decoded.Claims)
		return uuid.Nil, errors.New("could not parse token")
	}

	userIDValue, ok := claims["sub"]
	if !ok {
		return uuid.Nil, errors.New("userID not found in token payload")
	}
	userIDString, ok := userIDValue.(string)

	if !ok {
		return uuid.Nil, errors.New("could not parse ID")
	}

	userID, err := uuid.Parse(userIDString)

	if err != nil {
		return uuid.Nil, errors.New("could not parse ID")
	}

	return userID, nil
}
