package domain

import (
	"time"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type TicketStatus string

const (
	TicketStatusDraft       TicketStatus = "draft"
	TicketStatusOnSale      TicketStatus = "on_sale"
	TicketStatusPaused      TicketStatus = "paused"
	TicketStatusSoldOut     TicketStatus = "sold_out"
	TicketStatusUnavailable TicketStatus = "unavailable"
)

type Ticket struct {
	ID          uuid.UUID `json:"id"`
	EditionID   uuid.UUID `json:"edition_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`

	PriceCents         int  `json:"price_cents"`
	QuantitySold       int  `json:"quantity_sold"`
	QuantityReserved   int  `json:"quantity_reserved"`
	HasLimitedQuantity bool `json:"has_limited_quantity"`
	QuantityAvailable  int  `json:"quantity_available"`

	Status TicketStatus `json:"status"`

	CreatedBy uuid.UUID  `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateTicketSpec struct {
	EditionID   uuid.UUID `json:"edition_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`

	PriceCents         int  `json:"price_cents"`
	QuantityReserved   int  `json:"quantity_reserved"`
	HasLimitedQuantity bool `json:"has_limited_quantity"`
	QuantityAvailable  int  `json:"quantity_available"`

	MaxPerUser  *int `json:"max_per_user"`
	MinPerOrder int  `json:"min_per_order"`
}

func NewTicket(creatorID uuid.UUID, spec CreateTicketSpec) (*Ticket, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fail.New(errx.SYSUUIDV7GenerationError).WithArgs("NewTicket")
	}

	t := &Ticket{
		ID:                 id,
		EditionID:          spec.EditionID,
		Name:               spec.Name,
		Description:        spec.Description,
		PriceCents:         spec.PriceCents,
		QuantitySold:       0,
		QuantityReserved:   0,
		HasLimitedQuantity: spec.HasLimitedQuantity,
		QuantityAvailable:  spec.QuantityAvailable,
		Status:             TicketStatusDraft,
		CreatedBy:          creatorID,
	}

	if err := t.validate(); err != nil {
		return nil, err
	}

	return t, nil
}

// FIXME make me unique errors
func (t *Ticket) validate() error {
	if t.EditionID == uuid.Nil {
		return fail.New(errx.TicketValidationFailed).WithArgs("edition id is nil: " + uuid.Nil.String())
	}

	if t.Name == "" {
		return fail.New(errx.TicketValidationFailed).Trace("ticket name is required")
	}

	if t.HasLimitedQuantity && t.QuantityAvailable <= 0 {
		return fail.New(errx.TicketValidationFailed).Trace("cant have 0 or less available tickets at creation")
	}

	if t.PriceCents < 0 {
		return fail.New(errx.TicketValidationFailed).Trace("ticket price cannot be negative")
	}

	return nil
}
