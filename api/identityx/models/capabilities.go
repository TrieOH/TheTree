package models

import (
	"time"

	"github.com/google/uuid"
)

type Capability struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID *uuid.UUID `json:"project_id"`
	Resource  string     `json:"resource"`
	Action    string     `json:"action"`
	CreatedBy uuid.UUID  `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateCapabilityInput struct {
	Resource  string     `json:"resource"`
	Action    string     `json:"action"`
	ProjectID *uuid.UUID `json:"project_id"`
}
