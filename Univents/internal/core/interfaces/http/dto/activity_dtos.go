package dto

import (
	"time"
	"univents/internal/core/domain"

	"github.com/google/uuid"
)

type CreateActivityRequest struct {
	EditionScopeID uuid.UUID               `json:"edition_scope_id" validate:"required"`
	Title          string                  `json:"title" validate:"required,min=3"`
	Description    *string                 `json:"description"`
	Location       string                  `json:"location"`
	StartsAt       time.Time               `json:"starts_at" validate:"required"`
	EndsAt         time.Time               `json:"ends_at" validate:"required"`
	PresenterName  *string                 `json:"presenter_name"`
	TokenCost      int                     `json:"token_cost" validate:"gte=0"`
	HasCapacity    bool                    `json:"has_capacity"`
	Capacity       int                     `json:"capacity" validate:"gte=0"`
	Difficulty     *domain.DifficultyLevel `json:"difficulty"`
}
