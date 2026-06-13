package contracts

import (
	"encoding/json"
	"strings"
	"time"

	"payssage/internal/shared/errx"
	"payssage/internal/shared/validation"

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
	Provider    string          `json:"provider"`
	Metadata    json.RawMessage `json:"metadata"`

	// Provider Specific Data
	MercadoPagoData *MercadoPagoIntentData `json:"mercado_pago_data"`
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

type WebhookEndpoint struct {
	ID          uuid.UUID  `json:"id"`
	ScopeID     uuid.UUID  `json:"scope_id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	URL         string     `json:"url"`
	Secret      string     `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func NewWebhookEndpoint(workspaceID uuid.UUID, url, secret string) (*WebhookEndpoint, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("product").SetMessage("error generating uuid").SetCause(err)
	}

	w := &WebhookEndpoint{
		ID:          id,
		WorkspaceID: workspaceID,
		URL:         url,
		Secret:      secret,
		CreatedAt:   time.Now(),
	}

	if err := w.validate(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WebhookEndpoint) validate() error {
	return validation.Run(
		validation.RequireUUID("webhook_endpoint", "workspace_id", w.WorkspaceID),
		validation.RequireString("webhook_endpoint", "url", w.URL),
		validation.RequireString("webhook_endpoint", "secret", w.Secret),
		validation.Assert("webhook_endpoint", strings.HasPrefix(w.URL, "http://") || strings.HasPrefix(w.URL, "https://"), "url must be a valid http/https URL"),
	)
}

func (w *WebhookEndpoint) AddScope(scopeID uuid.UUID) {
	w.ScopeID = scopeID
}

type WebhookEventOriginal struct {
	ID          uuid.UUID
	WorkspaceID *uuid.UUID
	IntentID    *uuid.UUID
	Provider    string
	ExternalID  *string
	EventType   string
	Payload     []byte
	ReceivedAt  time.Time
}
