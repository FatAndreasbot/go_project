package models_test

import (
	"testing"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

func TestPasswordCheck(t *testing.T) {
	var nilUser *models.User = nil

	err := nilUser.CheckPassword("a password")
	if err == nil {
		t.Fatal("expected an error")
		return
	}
	// t.Log("asdasdasd")

	const strongPassword string = "a_very_strong_password_228"

	user := models.User{
		ID:      uuid.New(),
		Email:   "test@example.com",
		GroupID: uuid.New(),
	}

	user.SetPassword(strongPassword)
	err = user.CheckPassword("a password")
	if err == nil {
		t.Fatal("expected an error")
		return
	}
	t.Log(err.Error())

	err = user.CheckPassword(strongPassword)
	if err != nil {
		t.Fatal(err.Error())
		return
	}

}
