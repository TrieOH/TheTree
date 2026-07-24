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

type CreateTicketTypeRequest struct {
	Name        string  `json:"name"         validate:"required,min=2"`
	Description *string `json:"description"`
	AccessLevel int     `json:"access_level" validate:"gte=0"`
	PriceCents  int64   `json:"price_cents"  validate:"gte=0"`
	MaxQuantity *int    `json:"max_quantity" validate:"omitempty,gt=0"`
}

func (r CreateTicketTypeRequest) ToInput(editionID uuid.UUID) CreateTicketTypeInput {
	return CreateTicketTypeInput{
		EditionID:   editionID,
		Name:        r.Name,
		Description: r.Description,
		AccessLevel: r.AccessLevel,
		PriceCents:  r.PriceCents,
		MaxQuantity: r.MaxQuantity,
	}
}

type CreateTicketTypeInput struct {
	EditionID   uuid.UUID
	Name        string
	Description *string
	AccessLevel int
	PriceCents  int64
	MaxQuantity *int
}

type PatchTicketTypeRequest struct {
	Name        string  `json:"name"         validate:"required,min=2"`
	Description *string `json:"description"`
	AccessLevel int     `json:"access_level" validate:"gte=0"`
	PriceCents  int64   `json:"price_cents"  validate:"gte=0"`
	MaxQuantity *int    `json:"max_quantity" validate:"omitempty,gt=0"`
}

func (r PatchTicketTypeRequest) ToInput(ticketTypeID uuid.UUID) PatchTicketTypeInput {
	return PatchTicketTypeInput{
		TicketTypeID: ticketTypeID,
		Name:         r.Name,
		Description:  r.Description,
		AccessLevel:  r.AccessLevel,
		PriceCents:   r.PriceCents,
		MaxQuantity:  r.MaxQuantity,
	}
}

type PatchTicketTypeInput struct {
	TicketTypeID uuid.UUID
	Name         string
	Description  *string
	AccessLevel  int
	PriceCents   int64
	MaxQuantity  *int
}
