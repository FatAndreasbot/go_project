package grpcadapter

import (
	"context"
	"errors"

	"github.com/FatAndreasbot/go_project/user_service/infra/config"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
)

var publicMethods map[string]struct{} = map[string]struct{}{
	"/proto.user_service.UserService/LogIn": {},
}

const userdataKey = "userdata-e65f9095-b3d6-492f-86f4-7801c78c7732"

func Authenticate(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)
	if _, ok := publicMethods[method]; ok {
		return ctx, nil
	}

	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, errors.Join(err, errors.New("could not find jwt"))
	}

	// TODO
	userID, err := jwt.Parse(
		token,
		func(token *jwt.Token) (any, error) {
			return config.GetConfig().HS256_SECRET, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return ctx, errors.Join(err, errors.New("could not decode token"))
	}

	ctx = context.WithValue(ctx, userdataKey, userID)

	return ctx, nil
}
