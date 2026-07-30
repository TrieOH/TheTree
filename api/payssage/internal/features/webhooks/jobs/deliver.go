package jobs

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"payssage/ports"
	"time"

	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
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

// maxDeliveryAttempts is OUR retry budget for a delivery, tracked in
// webhook_deliveries.attempts. This is distinct from River's own
// per-job attempt tracking (InsertOpts.MaxAttempts) — River governs
// "should this job run again", this governs "have we given up on ever
// reaching this customer endpoint."
const maxDeliveryAttempts = 5

type DeliverWebhookWorker struct {
	river.WorkerDefaults[DeliverWebhookArgs]

	deliveries ports.WebhookDeliveryRepo
	events     ports.WebhookEventRepo
	endpoints  ports.WebhookEndpointRepo
}

func NewDeliverWebhookWorker(
	deliveries ports.WebhookDeliveryRepo,
	events ports.WebhookEventRepo,
	endpoints ports.WebhookEndpointRepo,
) *DeliverWebhookWorker {
	return &DeliverWebhookWorker{
		deliveries: deliveries,
		events:     events,
		endpoints:  endpoints,
	}
}

func (w *DeliverWebhookWorker) Work(ctx context.Context, job *river.Job[DeliverWebhookArgs]) error {
	delivery, err := w.deliveries.GetByID(ctx, job.Args.DeliveryID)
	if err != nil {
		telemetry.Log().Error("webhook delivery: delivery not found", zap.String("delivery_id", job.Args.DeliveryID.String()), zap.Error(err))
		return nil
	}

	if delivery.Status != models.WebhookDeliveryStatusPending {
		// Already delivered or already given up — a duplicate/stale job,
		// not an error.
		return nil
	}

	event, err := w.events.GetByID(ctx, delivery.EventID)
	if err != nil {
		return w.giveUp(ctx, delivery, 0, "", fmt.Sprintf("event not found: %v", err))
	}

	endpoint, err := w.endpoints.GetByID(ctx, delivery.EndpointID)
	if err != nil {
		return w.giveUp(ctx, delivery, 0, "", fmt.Sprintf("endpoint not found: %v", err))
	}

	signature := signPayload(endpoint.Secret, event.Payload)

	client := resty.New().SetTimeout(10 * time.Second)
	defer client.Close()

	resp, reqErr := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Payssage-Signature", signature).
		SetHeader("X-Payssage-Event-Type", event.EventType).
		SetBody([]byte(event.Payload)).
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
		telemetry.Log().Error("webhook delivery: failed to record success", zap.String("delivery_id", delivery.ID.String()), zap.Error(err))
		return err // genuinely worth retrying — we succeeded but couldn't record it
	}
	return nil
}

// recordAttempt records a failed attempt. If our own retry budget
// (maxDeliveryAttempts) isn't exhausted, it returns an error so River
// retries the job. Once exhausted, it marks the delivery permanently
// failed and returns nil so River stops.
func (w *DeliverWebhookWorker) recordAttempt(
	ctx context.Context,
	delivery *models.WebhookDelivery,
	attempts int,
	attemptedAt time.Time,
	responseStatus int,
	responseBody string,
) error {
	status := models.WebhookDeliveryStatusPending
	if attempts >= maxDeliveryAttempts {
		status = models.WebhookDeliveryStatusFailed
	}

	var respStatusPtr *int
	if responseStatus != 0 {
		respStatusPtr = &responseStatus
	}

	_, updateErr := w.deliveries.Update(ctx, models.UpdateDeliveryParams{
		ID:              delivery.ID,
		Status:          status,
		Attempts:        attempts,
		LastAttemptedAt: &attemptedAt,
		ResponseStatus:  respStatusPtr,
		ResponseBody:    &responseBody,
	})
	if updateErr != nil {
		telemetry.Log().Error("webhook delivery: failed to record attempt", zap.String("delivery_id", delivery.ID.String()), zap.Error(updateErr))
	}

	if status == models.WebhookDeliveryStatusFailed {
		telemetry.Log().Warn("webhook delivery: giving up after max attempts",
			zap.String("delivery_id", delivery.ID.String()),
			zap.Int("attempts", attempts),
		)
		return nil // stop River retries — we've made our own final decision
	}

	// Under our budget — return an error so River retries with its own
	// backoff policy.
	return fmt.Errorf("webhook delivery attempt %d failed: status %d", attempts, responseStatus)
}

// giveUp marks a delivery failed immediately, for errors that will never
// resolve by retrying (missing event/endpoint row).
func (w *DeliverWebhookWorker) giveUp(ctx context.Context, delivery *models.WebhookDelivery, responseStatus int, responseBody, reason string) error {
	telemetry.Log().Error("webhook delivery: giving up", zap.String("delivery_id", delivery.ID.String()), zap.String("reason", reason))
	_, _ = w.deliveries.Update(ctx, models.UpdateDeliveryParams{
		ID:              delivery.ID,
		Status:          models.WebhookDeliveryStatusFailed,
		Attempts:        delivery.Attempts + 1,
		LastAttemptedAt: new(time.Now()),
		ResponseBody:    &reason,
	})
	return nil
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
