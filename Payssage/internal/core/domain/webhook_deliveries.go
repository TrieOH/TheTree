package domain

import (
	"TriePayments/internal/shared/errx"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type DeliveryStatus string

const (
	DeliveryStatusPending   DeliveryStatus = "pending"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusFailed    DeliveryStatus = "failed"
)

const (
	EventPaymentSucceeded = "payment.succeeded"
	EventPaymentFailed    = "payment.failed"
	EventPaymentCancelled = "payment.cancelled"
)

type WebhookPayload struct {
	Event       string          `json:"event"`
	IntentID    uuid.UUID       `json:"intent_id"`
	WorkspaceID uuid.UUID       `json:"workspace_id"`
	Amount      int64           `json:"amount"`
	Currency    string          `json:"currency"`
	Metadata    json.RawMessage `json:"metadata"`
}

type WebhookDelivery struct {
	ID              uuid.UUID       `json:"id"`
	EndpointID      uuid.UUID       `json:"endpoint_id"`
	IntentID        uuid.UUID       `json:"intent_id"`
	Event           string          `json:"event"`
	Payload         json.RawMessage `json:"payload"`
	Status          DeliveryStatus  `json:"status"`
	Attempts        int             `json:"attempts"`
	LastAttemptedAt *time.Time      `json:"last_attempted_at"`
	CreatedAt       time.Time       `json:"created_at"`
}

func NewWebhookDelivery(endpointID, intentID uuid.UUID, event string, payload json.RawMessage) (*WebhookDelivery, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("product").SetMessage("error generating uuid").SetCause(err)
	}

	return &WebhookDelivery{
		ID:         id,
		EndpointID: endpointID,
		IntentID:   intentID,
		Event:      event,
		Payload:    payload,
		Status:     DeliveryStatusPending,
		CreatedAt:  time.Now(),
	}, nil
}

const TypeDeliverWebhook = "webhook:deliver"
const MaxDeliverRetries = 5

type DeliverWebhookPayload struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
	EndpointID uuid.UUID `json:"endpoint_id"`
	URL        string    `json:"url"`
	Secret     string    `json:"secret"`
	Payload    []byte    `json:"payload"`
}

func NewDeliverWebhookTask(deliveryID, endpointID uuid.UUID, url, secret string, payload []byte) (*asynq.Task, error) {
	data, err := json.Marshal(DeliverWebhookPayload{
		DeliveryID: deliveryID,
		EndpointID: endpointID,
		URL:        url,
		Secret:     secret,
		Payload:    payload,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeDeliverWebhook,
		data,
		asynq.MaxRetry(MaxDeliverRetries),
		asynq.TaskID(deliveryID.String()),
	), nil
}
