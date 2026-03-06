package dto

import (
	"encoding/json"
)

type EnsurePermissionsRequest struct {
	Permissions []PermissionDefinitionDTO `json:"permissions" validate:"required,min=1"`
}

type PermissionDefinitionDTO struct {
	Object string           `json:"object" validate:"required"`
	Action string           `json:"action" validate:"required"`
	Meta   *json.RawMessage `json:"meta"`
}

type EnsurePermissionsResponse struct {
	Permissions []EnsurePermissionResultDTO `json:"permissions"`
}

type EnsurePermissionResultDTO struct {
	Object  string `json:"object"`
	Action  string `json:"action"`
	Created bool   `json:"created"`
}
