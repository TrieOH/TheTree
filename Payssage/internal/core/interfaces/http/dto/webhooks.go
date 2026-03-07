package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProviderWebhookRequest struct {
	IntentID string `json:"intent_id" validate:"required"`
	Event    string `json:"event"     validate:"required"`
}

type RegisterWebhookEndpointRequest struct {
	URL string `json:"url" validate:"required"`
}

type WebhookEndpointResponse struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	URL         string    `json:"url"`
	Secret      string    `json:"secret"` // only shown on creation
	CreatedAt   time.Time `json:"created_at"`
}

type WebhookEndpointListResponse struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}
