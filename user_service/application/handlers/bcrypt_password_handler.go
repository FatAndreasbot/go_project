package handlers

import (
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordEncrypter struct{}

func (enc BcryptPasswordEncrypter) EncryptPassword(password string) (string, error) {
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12) // bcrypt.DefaultCost = 10
	if err != nil {
		return "", err
	}
	return string(newPasswordHash), nil
}

func (enc BcryptPasswordEncrypter) CheckPassword(password string, user *models.User) error {
	if user == nil {
		bcrypt.CompareHashAndPassword([]byte("definetly_a_valid_hash"), []byte(password))
		return bcrypt.ErrMismatchedHashAndPassword
	}
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

func (enc BcryptPasswordEncrypter) FakeCheckPassword(password string) {
	bcrypt.CompareHashAndPassword([]byte("definetly_a_valid_hash"), []byte(password))
}
