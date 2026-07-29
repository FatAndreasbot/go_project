package models

import "github.com/google/uuid"

type Group struct {
	ID          uuid.UUID
	Name        string
	Permossions []Permission
	Users       []User
}
