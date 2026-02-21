package dto

import "encoding/json"

type CreateScopeRequest struct {
	Name       string           `json:"name"`
	ExternalID *string          `json:"external_id"`
	ParentID   *string          `json:"parent_id"` // Optional: defaults to project root
	Meta       *json.RawMessage `json:"meta"`
}

type UpdateProjectScopeMetaRequest struct {
	Meta *json.RawMessage `json:"meta"`
}
