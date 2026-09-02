package grpcadapter

import (
	"errors"
	"log"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func EncodeAccessJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": config.GetConfig().AccessJWTExpiration,
		"iat": time.Now(),
	})

	signedToken, err := token.SignedString(config.GetConfig().HS256_SECRET)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func EncodeRefreshJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),

		"exp": config.GetConfig().RefreshJWTExpiration,
		"iat": time.Now(),
	})

	signedToken, err := token.SignedString(config.GetConfig().HS256_SECRET)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func DecodeAccessJWT(signedToken string) (uuid.UUID, error) {
	decoded, err := jwt.Parse(
		signedToken,
		func(token *jwt.Token) (any, error) {
			return config.GetConfig().HS256_SECRET, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)
	if !ok {
		log.Default().Printf("could not cast %v into MapClaims", decoded.Claims)
		return uuid.Nil, errors.New("could not parse token")
	}

	userIDValue, ok := claims["userID"]
	if !ok {
		return uuid.Nil, errors.New("userID not found in token payload")
	}
	userID := userIDValue.(uuid.UUID)

	return userID, errors.New("not implemented")
}
