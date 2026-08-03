package grpcadapter

import (
	"context"
	"errors"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
)

var publicMethods map[string]struct{} = map[string]struct{}{
	"/proto.user_service.UserService/LogIn": {},
}

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
	userID, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return nil, errors.New("not implemented")

	})
	if err != nil {
		return ctx, errors.Join(err, errors.New("could not decode token"))
	}

	ctx = context.WithValue(ctx, "userdata", userID)

	return ctx, nil
}
