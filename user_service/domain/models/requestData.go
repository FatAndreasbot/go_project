package models

import "github.com/google/uuid"

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type LoginResponseData struct {
	Tokens
	UserID uuid.UUID
}

type UserData struct {
	Email    string
	Password string
	GroupID  uuid.UUID
}

type UserUpdateRequestData struct {
	UserToChangeID uuid.UUID

	NewPassword string
	NewEmail    string
	NewGroupID  uuid.UUID
}
