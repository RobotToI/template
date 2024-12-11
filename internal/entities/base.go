package entities

import (
	"time"

	uuid "github.com/google/uuid"
)

// TimeModel is common entity part for CreatedAt and UpdatedAt fields
type TimeModel struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// IDModel is common entity part for ID(uint64 - Serial) fields
type IDModel struct {
	ID uint64 `db:"id"`
}

// UUIDModel is common entity part for UUID fields
type UUIDModel struct {
	ID uuid.UUID
}
