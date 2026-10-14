package grpcadapter_test

import (
	"log"
	"testing"

	grpcadapter "github.com/FatAndreasbot/go_project/user_service/adapters/inbound/grpc_adapter"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/config"
	"github.com/google/uuid"
)

func prepareTestData() *models.User {
	if err := config.LoadConfig("../../../.env"); err != nil {
		log.Default().Println(err.Error())
	}

	user := models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Group: &models.Group{
			ID:          uuid.New(),
			Name:        "test",
			Permissions: []*models.Permission{},
		},
	}

	return &user
}

func TestJWTEncoding(t *testing.T) {
	user := prepareTestData()

	jwt, err := grpcadapter.EncodeAccessJWT(user)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(jwt)
}

func TestJWTDecoding(t *testing.T) {
	user := prepareTestData()

	jwt, err := grpcadapter.EncodeAccessJWT(user)
	if err != nil {
		t.Error(err.Error())
	}

	userID, err := grpcadapter.DecodeJWT(jwt)
	if err != nil {
		t.Error(err.Error())
	}

	if userID != user.ID {
		t.Error("could not retrieve userID")
	}

}
