package jobs

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"lib/database"
	"log/slog"
	"payssage/ports"
	"time"
	"uuid"

	"lib/telemetry"
	"payssage/models"

	"resty.dev/v3"
)

// DeliverWebhookArgs carries only the delivery ID. The worker fetches
// everything else (delivery, event, endpoint) fresh at execution time —
// job args are not the place to carry secrets or payloads that may be
// stale by the time the job actually runs.
type DeliverWebhookArgs struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
}

func (DeliverWebhookArgs) Kind() string { return "webhook.deliver" }

// DeliverWebhookWorker POSTs one delivery to its tenant endpoint. Retries
// follow the jobs registry's policy, on every carrier: a failed attempt is
// recorded and returned as an error so the job is retried. When the retries
// run out, its give-up (webhook.fail) marks the delivery failed. The row is the source of
// truth, so a duplicate or stale job is a no-op.
type DeliverWebhookWorker struct {
	deliveries ports.WebhookDeliveryRepo
	events     ports.WebhookEventRepo
	endpoints  ports.WebhookEndpointRepo
	httpClient *resty.Client
}

func NewDeliverWebhookWorker(
	deliveries ports.WebhookDeliveryRepo,
	events ports.WebhookEventRepo,
	endpoints ports.WebhookEndpointRepo,
	httpClient *resty.Client,
) *DeliverWebhookWorker {
	return &DeliverWebhookWorker{
		deliveries: deliveries,
		events:     events,
		endpoints:  endpoints,
		httpClient: httpClient,
	}
}

// webhookEnvelope is the delivery body sent to tenant endpoints (D2): the
// correlation key is the Payssage intent id — never a provider-specific id.
// `payload` carries the raw provider event payload. The signature covers the
// envelope bytes, not the inner payload. `status_detail` is the normalized
// outcome detail (models.IntentStatusDetail, e.g. "high_risk") when the
// event failed/rejected — omitted otherwise.
type webhookEnvelope struct {
	IntentID     uuid.UUID       `json:"intent_id"`
	WalletID     uuid.UUID       `json:"wallet_id"`
	Provider     string          `json:"provider"`
	ExternalID   string          `json:"external_id"`
	EventType    string          `json:"event_type"`
	StatusDetail *string         `json:"status_detail,omitempty"`
	Payload      json.RawMessage `json:"payload"`
}

// buildEnvelope wraps a webhook event's raw provider payload in the tenant
// delivery envelope. The envelope is what gets signed and POSTed.
func buildEnvelope(event *models.WebhookEvent) ([]byte, error) {
	return json.Marshal(webhookEnvelope{
		IntentID:     event.IntentID,
		WalletID:     event.WalletID,
		Provider:     event.Provider,
		ExternalID:   event.ExternalID,
		EventType:    event.EventType,
		StatusDetail: event.StatusDetail,
		Payload:      event.Payload,
	})
}

func (w *DeliverWebhookWorker) Work(ctx context.Context, args DeliverWebhookArgs) error {
	delivery, err := w.deliveries.GetByID(ctx, args.DeliveryID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			telemetry.Log().ErrorContext(ctx, "webhook delivery: delivery not found", slog.String("delivery_id", args.DeliveryID.String()), telemetry.Err(err))
			return nil
		}
		return err // the database is unreachable — retry
	}

	if delivery.Status != models.WebhookDeliveryStatusPending {
		// Already delivered or already given up — a duplicate/stale job,
		// not an error.
		return nil
	}

	event, err := w.events.GetByID(ctx, delivery.EventID)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return err // the database is unreachable — retry
		}
		w.giveUp(ctx, delivery, 0, "", fmt.Sprintf("event not found: %v", err))
		return nil
	}

	endpoint, err := w.endpoints.GetByID(ctx, delivery.EndpointID)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			return err // the database is unreachable — retry
		}
		w.giveUp(ctx, delivery, 0, "", fmt.Sprintf("endpoint not found: %v", err))
		return nil
	}

	envelope, err := buildEnvelope(event)
	if err != nil {
		w.giveUp(ctx, delivery, 0, "", fmt.Sprintf("build envelope: %v", err))
		return nil
	}

	signature := signPayload(endpoint.Secret, envelope)

	resp, reqErr := w.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Payssage-Signature", signature).
		SetHeader("X-Payssage-Event-Type", event.EventType).
		SetBody(envelope).
		Post(endpoint.URL)

	attempts := delivery.Attempts + 1
	now := time.Now()

	if reqErr != nil {
		return w.recordAttempt(ctx, delivery, attempts, now, 0, reqErr.Error())
	}

	responseBody := truncateResponseBody(resp.String())

	if resp.IsStatusSuccess() {
		return w.recordSuccess(ctx, delivery, attempts, now, resp.StatusCode(), responseBody)
	}

	return w.recordAttempt(ctx, delivery, attempts, now, resp.StatusCode(), responseBody)
}

// recordSuccess marks a delivery as delivered.
func (w *DeliverWebhookWorker) recordSuccess(
	ctx context.Context,
	delivery *models.WebhookDelivery,
	attempts int,
	attemptedAt time.Time,
	responseStatus int,
	responseBody string,
) error {
	_, err := w.deliveries.Update(ctx, models.UpdateDeliveryParams{
		ID:              delivery.ID,
		Status:          models.WebhookDeliveryStatusDelivered,
		Attempts:        attempts,
		LastAttemptedAt: &attemptedAt,
		ResponseStatus:  &responseStatus,
		ResponseBody:    &responseBody,
	})
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "webhook delivery: failed to record success", slog.String("delivery_id", delivery.ID.String()), telemetry.Err(err))
		return err // genuinely worth retrying — we succeeded but couldn't record it
	}
	return nil
}

// recordAttempt records a failed attempt and returns an error, so the job
// is retried. The delivery stays pending: when the retries run out, the
// webhook.fail give-up marks it failed.
func (w *DeliverWebhookWorker) recordAttempt(
	ctx context.Context,
	delivery *models.WebhookDelivery,
	attempts int,
	attemptedAt time.Time,
	responseStatus int,
	responseBody string,
) error {
	var respStatusPtr *int
	if responseStatus != 0 {
		respStatusPtr = &responseStatus
	}

	_, updateErr := w.deliveries.Update(ctx, models.UpdateDeliveryParams{
		ID:              delivery.ID,
		Status:          models.WebhookDeliveryStatusPending,
		Attempts:        attempts,
		LastAttemptedAt: &attemptedAt,
		ResponseStatus:  respStatusPtr,
		ResponseBody:    &responseBody,
	})
	if updateErr != nil {
		telemetry.Log().ErrorContext(ctx, "webhook delivery: failed to record attempt", slog.String("delivery_id", delivery.ID.String()), telemetry.Err(updateErr))
	}

	return fmt.Errorf("webhook delivery attempt %d failed: status %d", attempts, responseStatus)
}

// giveUp marks a delivery failed immediately, for errors that will never
// resolve by retrying (missing event/endpoint row).
func (w *DeliverWebhookWorker) giveUp(ctx context.Context, delivery *models.WebhookDelivery, responseStatus int, responseBody, reason string) {
	telemetry.Log().ErrorContext(ctx, "webhook delivery: giving up", slog.String("delivery_id", delivery.ID.String()), slog.String("reason", reason))
	_, _ = w.deliveries.Update(ctx, models.UpdateDeliveryParams{
		ID:              delivery.ID,
		Status:          models.WebhookDeliveryStatusFailed,
		Attempts:        delivery.Attempts + 1,
		LastAttemptedAt: new(time.Now()),
		ResponseBody:    &responseBody,
		ResponseStatus:  &responseStatus,
	})
}

func signPayload(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func truncateResponseBody(body string) string {
	const maxLen = 10000
	if len(body) > maxLen {
		return body[:maxLen]
	}
	return body
}
