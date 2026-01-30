package dto

import "github.com/google/uuid"

type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

type UpdateRoleRequest struct {
	Description *string `json:"description"`
}

type UserRoleRequest struct {
	RoleID  uuid.UUID  `json:"role_id" validate:"required"`
	ScopeID *uuid.UUID `json:"scope_id"`
}
