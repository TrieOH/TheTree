package jobs

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"lib/database"
	"lib/telemetry"
	"payssage/models"
	"payssage/ports"
)

// FailWebhookDeliveryArgs is sent when the carrier has given up on a
// webhook.deliver job: on Cloudflare, the Worker's dead-letter consumer
// turns each dead webhook.deliver message into one of these.
type FailWebhookDeliveryArgs struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
}

func (FailWebhookDeliveryArgs) Kind() string { return "webhook.fail" }

// FailWebhookDeliveryWorker marks a still-pending delivery failed, keeping
// the last attempt's response on the row.
type FailWebhookDeliveryWorker struct {
	deliveries ports.WebhookDeliveryRepo
}

func NewFailWebhookDeliveryWorker(deliveries ports.WebhookDeliveryRepo) *FailWebhookDeliveryWorker {
	return &FailWebhookDeliveryWorker{deliveries: deliveries}
}

func (w *FailWebhookDeliveryWorker) Work(ctx context.Context, args FailWebhookDeliveryArgs) error {
	delivery, err := w.deliveries.GetByID(ctx, args.DeliveryID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil
		}
		return err
	}
	if delivery.Status != models.WebhookDeliveryStatusPending {
		return nil
	}

	telemetry.Log().WarnContext(ctx, "webhook delivery: giving up after retries",
		slog.String("delivery_id", delivery.ID.String()),
		slog.Int("attempts", delivery.Attempts),
	)
	_, err = w.deliveries.Update(ctx, models.UpdateDeliveryParams{
		ID:              delivery.ID,
		Status:          models.WebhookDeliveryStatusFailed,
		Attempts:        delivery.Attempts,
		LastAttemptedAt: delivery.LastAttemptedAt,
		ResponseStatus:  delivery.ResponseStatus,
		ResponseBody:    delivery.ResponseBody,
	})
	return err
}
