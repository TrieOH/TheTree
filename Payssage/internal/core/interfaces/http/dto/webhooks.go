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

type MercadoPagoWebhookRequest struct {
	Action string `json:"action"`
	Data   struct {
		ID string `json:"id"`
	} `json:"data"`
}

type WebhookDeliveryResponse struct {
	ID              uuid.UUID  `json:"id"`
	EndpointID      uuid.UUID  `json:"endpoint_id"`
	IntentID        uuid.UUID  `json:"intent_id"`
	Event           string     `json:"event"`
	Status          string     `json:"status"`
	Attempts        int        `json:"attempts"`
	LastAttemptedAt *time.Time `json:"last_attempted_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type WebhookEventResponse struct {
	ID          uuid.UUID  `json:"id"`
	Provider    string     `json:"provider"`
	EventType   string     `json:"event_type"`
	ExternalID  *string    `json:"external_id,omitempty"`
	WorkspaceID *uuid.UUID `json:"workspace_id,omitempty"`
	IntentID    *uuid.UUID `json:"intent_id,omitempty"`
	ReceivedAt  time.Time  `json:"received_at"`
}
