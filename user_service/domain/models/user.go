package models

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID
	Name         string
	passwordHash string
	Group        Group
}

func NewUser(name, password string, groupID uuid.UUID) (*User, error) {
	user := User{
		Name: name,
	}

	return &user, errors.New("not implemented")
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
}
