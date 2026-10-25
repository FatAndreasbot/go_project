package models

import (
	"time"

	"github.com/google/uuid"
)

type Cursor struct {
	Timestamp time.Time
	ID        uuid.UUID
}
