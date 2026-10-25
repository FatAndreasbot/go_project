package grpcadapter

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/application/usecases"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
)

var publicMethods []string = []string{
	"/user_service.v1.UserService/LogIn",
	"/user_service.v1.UserService/Refresh",
}

var permissionSettings map[string]models.PermissionName = map[string]models.PermissionName{}

type contextValueKey string

// points to userID: uuid.UUID
const userdataKey contextValueKey = "userdata-e65f9095-b3d6-492f-86f4-7801c78c7732"

func GetAuthInterceptor(gen *usecases.MiddlewareGenerator) grpc.UnaryServerInterceptor {
	return auth.UnaryServerInterceptor(
		auth.AuthFunc(gen.GetAuthFunction(publicMethods)),
	)
}

func GetPermissionCheckingInterceptor(gen *usecases.MiddlewareGenerator) grpc.UnaryServerInterceptor {
	return auth.UnaryServerInterceptor(
		auth.AuthFunc(gen.GetPermissionFunction(permissionSettings)),
	)
}

func GetEndpoint(ctx context.Context) (string, error) {
	method, _ := grpc.Method(ctx)
	return method, nil
}

// type AuthInterceptor struct {
// 	jwtHandler    application.JWTHandlerPort
// 	sessionHander outgoing.SessionHandler
// }

// func NewAuthInterceptor(
// 	jwtHandler application.JWTHandlerPort,
// 	sessiohadler outgoing.SessionHandler,
// ) *AuthInterceptor {
// 	return &AuthInterceptor{
// 		jwtHandler:    jwtHandler,
// 		sessionHander: sessiohadler,
// 	}
// }

// func (factory *AuthInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
// 	return auth.UnaryServerInterceptor(
// 		func(ctx context.Context) (context.Context, error) {
// 			jwtHandler := factory.jwtHandler
// 			sessionHandler := factory.sessionHander

// 			method, _ := grpc.Method(ctx)
// 			if _, ok := publicMethods[method]; ok {
// 				return ctx, nil
// 			}

// 			token, err := auth.AuthFromMD(ctx, "bearer")
// 			if err != nil {
// 				return ctx, errors.Join(err, errors.New("could not find jwt"))
// 			}

// 			userID, err := jwtHandler.DecodeJWT(ctx, token)
// 			if err != nil {
// 				log.Default().Println(err)
// 				return ctx, errors.New("error when reading parsing userID")
// 			}

// 			if err := sessionHandler.CheckSession(ctx, userID); err != nil {
// 				return ctx, err
// 			}

// 			ctx = context.WithValue(ctx, userdataKey, userID)
// 			return ctx, nil
// 		},
// 	)
// }
