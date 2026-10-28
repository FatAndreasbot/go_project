package grpcadapter

import (
	"context"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"
)

var PublicMethods []string = []string{
	"/user_service.v1.UserService/LogIn",
	"/user_service.v1.UserService/Refresh",
}

var PermissionSettings map[string]models.PermissionName = map[string]models.PermissionName{
	// Login stuff
	"/user_service.v1.UserService/LogIn":   models.Any,
	"/user_service.v1.UserService/LogOut":  models.Any,
	"/user_service.v1.UserService/Refresh": models.Any,
	// User CRUD
	"/user_service.v1.UserService/CreateUser": models.UserCreate,
	"/user_service.v1.UserService/GetUser":    models.UserRead,
	"/user_service.v1.UserService/GetUsers":   models.UserRead,
	"/user_service.v1.UserService/UpdateUser": models.UserUpdate,
	"/user_service.v1.UserService/DeleteUser": models.UserDelete,
	// Group CRUD
	"/user_service.v1.UserService/CreateGroup":               models.GroupCreate,
	"/user_service.v1.UserService/GetGroup":                  models.GroupRead,
	"/user_service.v1.UserService/GetGroups":                 models.GroupRead,
	"/user_service.v1.UserService/UpdateGroup":               models.GroupUpdate,
	"/user_service.v1.UserService/AddPermissionToGroup":      models.GroupUpdate,
	"/user_service.v1.UserService/RemovePermissionFromGroup": models.GroupUpdate,
	"/user_service.v1.UserService/DeleteGroup":               models.GroupDelete,
	// Permissions
	"/user_service.v1.UserService/GetAllPermissions": models.Any,
}

type contextValueKey string

// points to userID: uuid.UUID
// const userdataKey contextValueKey = "userdata-e65f9095-b3d6-492f-86f4-7801c78c7732"

// func GetAuthInterceptor(gen *usecases.MiddlewareGenerator) grpc.UnaryServerInterceptor {
// 	return auth.UnaryServerInterceptor(
// 		auth.AuthFunc(gen.GetAuthFunction(PublicMethods)),
// 	)
// }

// func GetPermissionCheckingInterceptor(gen *usecases.MiddlewareGenerator) grpc.UnaryServerInterceptor {
// 	return auth.UnaryServerInterceptor(
// 		auth.AuthFunc(gen.GetPermissionFunction(PermissionSettings)),
// 	)
// }

func GetEndpoint(ctx context.Context) (string, error) {
	method, _ := grpc.Method(ctx)
	return method, nil
}

func GetToken(ctx context.Context) (string, error) {
	return auth.AuthFromMD(ctx, "bearer")
}
