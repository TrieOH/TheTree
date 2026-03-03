package domain

import (
	"time"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type Ticket struct {
	ID          uuid.UUID `json:"id"`
	EditionID   uuid.UUID `json:"edition_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`

	CreatedBy uuid.UUID  `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateTicketSpec struct {
	EditionScopeID uuid.UUID `json:"edition_scope_id"`
	EditionID      uuid.UUID `json:"edition_id"`
	Name           string    `json:"name"`
	Description    *string   `json:"description"`
}

func NewTicket(creatorID uuid.UUID, spec CreateTicketSpec) (*Ticket, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fail.New(errx.SYSUUIDV7GenerationError).WithArgs("NewTicket")
	}

	t := &Ticket{
		ID:          id,
		EditionID:   spec.EditionID,
		Name:        spec.Name,
		Description: spec.Description,
		CreatedBy:   creatorID,
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

	if t.CreatedBy == uuid.Nil {
		return fail.New(errx.TicketValidationFailed).WithArgs("created by id is nil: " + uuid.Nil.String())
	}

	if t.Name == "" {
		return fail.New(errx.TicketValidationFailed).Trace("ticket name is required")
	}

	return nil
}
