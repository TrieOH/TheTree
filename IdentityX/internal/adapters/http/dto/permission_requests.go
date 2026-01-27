package dto

import "encoding/json"

type CreatePermissionRequest struct {
	Object     string           `json:"object" validate:"required"`
	Action     string           `json:"action" validate:"required"`
	Conditions *json.RawMessage `json:"conditions"`
}
