package dto

import (
	"time"
	"univents/internal/core/domain"
)

type CreateCheckpointRequest struct {
	StartsAt   *time.Time              `json:"starts_at"`
	EndsAt     *time.Time              `json:"ends_at"`
	Name       string                  `json:"name"`
	Type       domain.CheckpointType   `json:"type"`
	AccessMode domain.CheckpointAccess `json:"access_mode" validate:"required"`
}
