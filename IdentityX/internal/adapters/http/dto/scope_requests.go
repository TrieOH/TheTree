package dto

type CreateScopeRequest struct {
	Name       string  `json:"name"`
	ExternalID *string `json:"external_id"`
	ParentID   *string `json:"parent_id"` // Optional: defaults to project root
}
