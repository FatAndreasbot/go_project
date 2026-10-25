package models

import (
	"time"

	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken      string
	AccesExpiresAt   time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type LoginResponseData struct {
	Tokens
	UserID uuid.UUID
}

type UserData struct {
	Email    string
	Password string
	GroupID  uuid.UUID

	IdempotencyKey uuid.UUID
}

type UserUpdateRequestData struct {
	UserToChangeID uuid.UUID

	NewPassword string
	NewEmail    string
	NewGroupID  uuid.UUID
}
