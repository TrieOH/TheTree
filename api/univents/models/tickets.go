package models

import (
	"time"

	"github.com/google/uuid"
)

type TicketType struct {
	ID          uuid.UUID  `json:"id"`
	EditionID   uuid.UUID  `json:"edition_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	AccessLevel int        `json:"access_level"`
	PriceCents  int64      `json:"price_cents"`
	MaxQuantity *int       `json:"max_quantity"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type CreateTicketTypeInput struct {
	EditionID   uuid.UUID
	Name        string
	Description *string
	AccessLevel int
	PriceCents  int64
	MaxQuantity *int
}

type PatchTicketTypeInput struct {
	TicketTypeID uuid.UUID
	Name         string
	Description  *string
	AccessLevel  int
	PriceCents   int64
	MaxQuantity  *int
}
