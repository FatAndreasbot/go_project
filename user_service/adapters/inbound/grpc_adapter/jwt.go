package grpcadapter

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func EncodeAccessJWT(user *models.User) (string, error) {
	tokenIssuedAt := time.Now()

	expTimeString, err := config.GetConfig(config.AccessJWTExpiration)
	if err != nil {
		return "", err
	}

	expTimeSeconds, err := strconv.Atoi(expTimeString)
	if err != nil {
		return "", err
	}
	expTime := time.Second * time.Duration(expTimeSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(tokenIssuedAt.Add(expTime)),
		IssuedAt:  jwt.NewNumericDate(tokenIssuedAt),
	})

	secret, err := config.GetConfig(config.HS256_SECRET)
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func EncodeRefreshJWT(user *models.User) (string, error) {
	tokenIssuedAt := time.Now()

	expTimeString, err := config.GetConfig(config.RefreshJWTExpiration)
	if err != nil {
		return "", err
	}

	expTimeSeconds, err := strconv.Atoi(expTimeString)
	if err != nil {
		return "", err
	}
	expTime := time.Second * time.Duration(expTimeSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(tokenIssuedAt.Add(expTime)),
		IssuedAt:  jwt.NewNumericDate(tokenIssuedAt),
	})

	secret, err := config.GetConfig(config.HS256_SECRET)
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func DecodeJWT(signedToken string) (uuid.UUID, error) {
	secret, err := config.GetConfig(config.HS256_SECRET)
	if err != nil {
		return uuid.Nil, err
	}

	decoded, err := jwt.Parse(
		signedToken,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
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
