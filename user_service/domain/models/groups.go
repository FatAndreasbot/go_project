package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID          uuid.UUID
	Name        string
	Permissions []*Permission

	CreatedAt time.Time
}
