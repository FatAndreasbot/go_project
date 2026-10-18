package models

import (
	"time"

	"github.com/google/uuid"
)

type MainGroupData struct {
	ID   uuid.UUID
	Name string

	CreatedAt time.Time
}

type Group struct {
	MainGroupData

	Permissions []*Permission
}
