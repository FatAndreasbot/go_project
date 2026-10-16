package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	GroupID      uuid.UUID
	CreatedAt    time.Time
}

func (u *User) CheckPassword(password string) error {
	if u == nil {
		bcrypt.CompareHashAndPassword([]byte("definetly_a_valid_hash"), []byte(password))
		return bcrypt.ErrMismatchedHashAndPassword
	}
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
