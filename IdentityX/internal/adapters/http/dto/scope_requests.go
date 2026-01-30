package dto

type CreateScopeRequest struct {
	Name       string  `json:"name"`
	ExternalID *string `json:"external_id"`
}
