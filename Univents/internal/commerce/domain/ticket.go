package domain

import (
	"time"
	"univents/internal/shared/errx"
	"univents/internal/shared/validation"

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
		return nil, errx.Internal("ticket").SetMessage("error generating uuid").SetCause(err)
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

func (t *Ticket) validate() error {
	return validation.Run(
		validation.RequireUUID("ticket", "edition_id", t.EditionID),
		validation.RequireUUID("ticket", "created_by", t.CreatedBy),
		validation.RequireString("ticket", "name", t.Name),
	)
}
