package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WebhookEndpoint struct {
	ID        uuid.UUID `json:"id"`
	WalletID  uuid.UUID `json:"wallet_id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Secret    string    `json:"secret"`
	CreatedAt time.Time `json:"created_at"`
}

type WebhookEvent struct {
	ID         uuid.UUID       `json:"id"`
	WalletID   uuid.UUID       `json:"wallet_id"`
	IntentID   uuid.UUID       `json:"intent_id"`
	Provider   string          `json:"provider"`
	ExternalID string          `json:"external_id"`
	EventType  string          `json:"event_type"`
	Payload    json.RawMessage `json:"payload"`
	ReceivedAt time.Time       `json:"received_at"`
}

type WebhookDeliveryStatus string

const (
	WebhookDeliveryStatusPending   WebhookDeliveryStatus = "pending"
	WebhookDeliveryStatusDelivered WebhookDeliveryStatus = "delivered"
	WebhookDeliveryStatusFailed    WebhookDeliveryStatus = "failed"
)

type WebhookDeliveries struct {
	ID              uuid.UUID             `json:"id"`
	EndpointID      uuid.UUID             `json:"endpoint_id"`
	EventID         uuid.UUID             `json:"event_id"`
	Status          WebhookDeliveryStatus `json:"status"`
	Attempts        int                   `json:"attempts"`
	LastAttemptedAt *time.Time            `json:"last_attempted_at"`
	ResponseStatus  *int                  `json:"response_status"`
	ResponseBody    *string               `json:"response_body"`
	CreatedAt       time.Time             `json:"created_at"`
}
