package usecases

import (
	"context"
	"errors"
	"log"
	"slices"

	domainerrors "github.com/FatAndreasbot/go_project/user_service/domain/errors"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/application"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"
	"github.com/google/uuid"
)

type MiddlewareGenerator struct {
	jwtHandler        application.JWTHandlerPort
	sessionHander     outgoing.SessionHandler
	permissionHandler application.PermissionHandlerPort

	getEndpoint func(context.Context) (string, error)
	getToken    func(context.Context) (string, error)
}

func NewMiddlewareGenerator(
	jwtHandler application.JWTHandlerPort,
	sessionHander outgoing.SessionHandler,
	permissionHandler application.PermissionHandlerPort,
	getEndpoint func(context.Context) (string, error),
	getToken func(context.Context) (string, error),
) *MiddlewareGenerator {
	return &MiddlewareGenerator{
		jwtHandler:        jwtHandler,
		sessionHander:     sessionHander,
		permissionHandler: permissionHandler,
		getEndpoint:       getEndpoint,
		getToken:          getToken,
	}
}

func (mg *MiddlewareGenerator) GetAuthFunction(publicMethods []string) models.Interceptor {

	return func(ctx context.Context) (context.Context, error) {
		jwtHandler := mg.jwtHandler
		sessionHandler := mg.sessionHander

		method, err := mg.getEndpoint(ctx)
		if err != nil {
			return ctx, err
		}

		if slices.Contains(publicMethods, method) {
			return ctx, nil
		}

		token, err := mg.getToken(ctx)
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

		ctx = context.WithValue(ctx, models.UserdataKey, userID)
		return ctx, nil
	}
}

func (mg *MiddlewareGenerator) GetPermissionFunction(
	permissionSettigns map[string]models.PermissionName,
	publicMethods []string,
) models.Interceptor {
	return func(ctx context.Context) (context.Context, error) {
		method, err := mg.getEndpoint(ctx)
		if err != nil {
			return ctx, err
		}
		requesterID, ok := ctx.Value(models.UserdataKey).(uuid.UUID)
		if !ok {
			if slices.Contains(publicMethods, method) {
				return ctx, nil
			}
			log.Default().Printf("endoint %s has no associated permission\n", method)
			return ctx, errors.New("unautorized")
		}

		if !mg.permissionHandler.HasPermission(ctx, requesterID, permissionSettigns[method]) {
			return ctx, domainerrors.NotEnoughPermission
		}

		return ctx, nil
	}
}
