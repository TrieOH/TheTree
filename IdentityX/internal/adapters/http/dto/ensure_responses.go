package dto

import (
	"encoding/json"
)

type EnsureRolesRequest struct {
	Roles []RoleDefinitionDTO `json:"roles" validate:"required,min=1"`
}

type RoleDefinitionDTO struct {
	Name        string                    `json:"name" validate:"required"`
	Permissions []PermissionDefinitionDTO `json:"permissions"`
	Meta        *json.RawMessage          `json:"meta"`
}

type EnsureRolesResponse struct {
	Roles []EnsureRoleResultDTO `json:"roles"`
}

type EnsureRoleResultDTO struct {
	Name    string `json:"name"`
	Created bool   `json:"created"`
}
