package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreatePermissionRequest struct {
	Object     string           `json:"object" validate:"required"`
	Action     string           `json:"action" validate:"required"`
	Conditions *json.RawMessage `json:"conditions"`
}

type UserPermissionRequest struct {
	ScopeID      *uuid.UUID `json:"scope_id"`
	PermissionID uuid.UUID  `json:"permission_id" validate:"required"`
}

type CheckRequest struct {
	ProjectID *uuid.UUID `json:"project_id"`
	ScopeID   *uuid.UUID `json:"scope_id"`
	EntityID  uuid.UUID  `json:"entity_id" validate:"required"`
	Object    string     `json:"object" validate:"required"`
	Action    string     `json:"action" validate:"required"`
}
