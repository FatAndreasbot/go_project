package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID
	Name         string
	PasswordHash string
	Group        *Group
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

func (u *User) SetPassword(password string) error {
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12) // bcrypt.DefaultCost = 10
	if err != nil {
		return err
	}

	u.PasswordHash = string(newPasswordHash)
	return nil
}
