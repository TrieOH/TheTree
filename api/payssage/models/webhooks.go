package models

import (
	"encoding/json"
	"net/http"
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

type WebhookDelivery struct {
	ID              uuid.UUID             `json:"id"`
	EndpointID      uuid.UUID             `json:"endpoint_id"`
	EventID         uuid.UUID             `json:"event_id"`
	Status          WebhookDeliveryStatus `json:"status"`
	Attempts        int                   `json:"attempts"`
	LastAttemptedAt *time.Time            `json:"last_attempted_at"`
	ResponseStatus  *int                  `json:"response_status"`
	ResponseBody    *string               `json:"response_body"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       *time.Time            `json:"updated_at"`
}

type WebhookParseResult struct {
	WalletID   uuid.UUID
	IntentID   uuid.UUID
	EventType  string // normalized, e.g. "payment.succeeded" — not the provider's raw type string
	ExternalID string // provider's own event/resource id, used for the webhook_events.sql dedupe index
}

type ReceiveWebhookInput struct {
	Provider string
	Request  *http.Request
	RawBody  []byte
}

// UpdateDeliveryParams carries the fields written after a delivery
// attempt (success or failure). Status, Attempts, and LastAttemptedAt
// are always set; ResponseStatus/ResponseBody are pointers since a
// request-level failure (timeout, connection refused) never produces
// an HTTP response to record.
type UpdateDeliveryParams struct {
	ID              uuid.UUID
	Status          WebhookDeliveryStatus
	Attempts        int
	LastAttemptedAt *time.Time
	ResponseStatus  *int
	ResponseBody    *string
}

type CreateWebhookEndpointRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (r CreateWebhookEndpointRequest) ToInput(walletID uuid.UUID) CreateWebhookEndpointInput {
	return CreateWebhookEndpointInput{
		WalletID: walletID,
		Name:     r.Name,
		URL:      r.URL,
	}
}

type CreateWebhookEndpointInput struct {
	WalletID uuid.UUID
	Name     string
	URL      string
}
