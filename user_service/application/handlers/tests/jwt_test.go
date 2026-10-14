package handlers_test

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/application/handlers"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/config"
	"github.com/google/uuid"
)

func prepareTestData() (*models.User, *handlers.JWTHandler) {
	if err := config.LoadConfig("../../../.env"); err != nil {
		log.Default().Println(err.Error())
	}

	user := models.User{
		ID:      uuid.New(),
		Email:   "test@example.com",
		GroupID: uuid.New(),
	}

	return &user, handlers.NewJWTHandler()
}

func TestJWTEncoding(t *testing.T) {
	user, jwtHandler := prepareTestData()
	ctx := context.Background()

	jwt, err := jwtHandler.EncodeAccessJWT(ctx, user)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(jwt)
}

func TestJWTDecoding(t *testing.T) {
	user, jwtHandler := prepareTestData()
	ctx := context.Background()

	jwt, err := jwtHandler.EncodeAccessJWT(ctx, user)
	if err != nil {
		t.Error(err.Error())
	}

	userID, err := jwtHandler.DecodeJWT(ctx, jwt)
	if err != nil {
		t.Error(err.Error())
	}

	if userID != user.ID {
		t.Error("could not retrieve userID")
	}
}

func TestJWTTimeout(t *testing.T) {
	user, jwtHandler := prepareTestData()
	ctx := context.Background()

	jwt, err := jwtHandler.EncodeAccessJWT(ctx, user)
	if err != nil {
		t.Error(err.Error())
	}
	expTime, err := config.GetConfig(config.AccessJWTExpiration)
	if err != nil {
		t.Error(err.Error())
	}
	timeToSleepSeconds, err := strconv.Atoi(expTime)
	if err != nil {
		t.Error(err.Error())
	}

	t.Logf("sleeping for %d seconds", timeToSleepSeconds+1)

	time.Sleep(time.Second * time.Duration(timeToSleepSeconds+1))

	_, err = jwtHandler.DecodeJWT(ctx, jwt)
	if err == nil {
		t.Error("expected to fail")
	}
	t.Log(err.Error())
}
