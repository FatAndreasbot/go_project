package application

import "github.com/FatAndreasbot/go_project/user_service/domain/models"

type PasswordEncrypter interface {
	EncryptPassword(string) (string, error)
	CheckPassword(string, *models.User) error
	FakeCheckPassword(string)
}
