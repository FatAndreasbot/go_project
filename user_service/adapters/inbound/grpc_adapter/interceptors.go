package grpcadapter

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
)

var publicMethods map[string]struct{} = map[string]struct{}{
	"/proto.user_service.UserService/LogIn": {},
}

// points to userID: uuid.UUID
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

	userID, err := DecodeJWT(token)
	if err != nil {
		return ctx, errors.New("error when reading parsing userID")
	}

	ctx = context.WithValue(ctx, userdataKey, userID)

	return ctx, nil
}
