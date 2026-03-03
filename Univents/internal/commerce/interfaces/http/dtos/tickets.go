package dtos

import "github.com/google/uuid"

type CreateTicketRequest struct {
	EditionScopeID uuid.UUID `json:"edition_scope_id" validate:"required"`
	Name           string    `json:"name" validate:"required,min=3"`
	Description    *string   `json:"description"`
}
