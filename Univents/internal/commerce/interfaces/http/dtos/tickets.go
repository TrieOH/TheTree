package dtos

import (
	"univents/internal/commerce/domain"

	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	EditionScopeID uuid.UUID `json:"edition_scope_id" validate:"required"`
	Name           string    `json:"name" validate:"required,min=3"`
	Description    *string   `json:"description"`
}

type AddTicketPermissionRequest struct {
	TicketScopeID  uuid.UUID             `json:"ticket_scope_id" validate:"required"`
	PermissionType domain.PermissionType `json:"permission_type"`
	ActivityID     *uuid.UUID            `json:"activity_id"`
	ProductID      *uuid.UUID            `json:"product_id"`
	CheckpointID   *uuid.UUID            `json:"checkpoint_id"`
}

type RemoveTicketPermissionRequest struct {
	TicketScopeID uuid.UUID `json:"ticket_scope_id" validate:"required"`
}
