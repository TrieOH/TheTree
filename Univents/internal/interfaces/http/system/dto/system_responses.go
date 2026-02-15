package dto

import "github.com/google/uuid"

type HealthResponse struct {
	Status  string    `json:"status" example:"ok"`
	Service string    `json:"service" example:"univents-api"`
	UserID  uuid.UUID `json:"user_id,omitempty" example:"some-uuid"`
}
