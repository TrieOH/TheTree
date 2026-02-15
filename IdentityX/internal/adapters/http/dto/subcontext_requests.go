package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type AddSubContextRequest struct {
	UserID uuid.UUID       `json:"user_id" validate:"required"`
	Data   json.RawMessage `json:"data" validate:"required"`
}

type RemoveSubContextRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Keys   []string  `json:"keys" validate:"required,min=1"`
}
