package dto

import (
	"time"
	"univents/internal/core/domain"

	"github.com/google/uuid"
)

type CreateCheckpointRequest struct {
	EditionScopeID uuid.UUID               `json:"edition_scope_id"`
	StartsAt       *time.Time              `json:"starts_at"`
	EndsAt         *time.Time              `json:"ends_at"`
	Name           string                  `json:"name"`
	Type           domain.CheckpointType   `json:"type"`
	AccessMode     domain.CheckpointAccess `json:"access_mode" validate:"required"`
}
