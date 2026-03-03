package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	Name        string           `json:"name" validate:"required"`
	Description *string          `json:"description"`
	Meta        *json.RawMessage `json:"meta"`
}

type UpdateRoleDescriptionRequest struct {
	Description *string `json:"description"`
}

type UpdateRoleMetaRequest struct {
	Meta *json.RawMessage `json:"meta"`
}

type UserRoleRequest struct {
	RoleID  uuid.UUID  `json:"role_id" validate:"required"`
	ScopeID *uuid.UUID `json:"scope_id"`
}

type UserRoleByNameRequest struct {
	RoleName string     `json:"role_name" validate:"required"`
	ScopeID  *uuid.UUID `json:"scope_id"`
}
