package grpcadapter

import (
	"context"
	"errors"
	"log"

	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
)

var publicMethods map[string]struct{} = map[string]struct{}{
	"/user_service.v1.UserService/LogIn":   {},
	"/user_service.v1.UserService/Refresh": {},
}

// points to userID: uuid.UUID
const userdataKey = "userdata-e65f9095-b3d6-492f-86f4-7801c78c7732"

type AuthInterceptor struct {
	jwtHandler    application.JWTHandler
	sessionHander outgoing.SessionHandler
}

func NewAuthInterceptor(
	jwtHandler application.JWTHandler,
	sessiohadler outgoing.SessionHandler,
) *AuthInterceptor {
	return &AuthInterceptor{
		jwtHandler:    jwtHandler,
		sessionHander: sessiohadler,
	}
}

func (factory *AuthInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return auth.UnaryServerInterceptor(
		func(ctx context.Context) (context.Context, error) {
			jwtHandler := factory.jwtHandler
			sessionHandler := factory.sessionHander

			method, _ := grpc.Method(ctx)
			if _, ok := publicMethods[method]; ok {
				return ctx, nil
			}

			token, err := auth.AuthFromMD(ctx, "bearer")
			if err != nil {
				return ctx, errors.Join(err, errors.New("could not find jwt"))
			}

			userID, err := jwtHandler.DecodeJWT(ctx, token)
			if err != nil {
				log.Default().Println(err)
				return ctx, errors.New("error when reading parsing userID")
			}

			if err := sessionHandler.CheckSession(ctx, userID); err != nil {
				return ctx, err
			}

			ctx = context.WithValue(ctx, userdataKey, userID)
			return ctx, nil
		},
	)
}
