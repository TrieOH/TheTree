package jobs

import (
	"context"
	"errors"
	"log/slog"
	"uuid"

	"lib/database"
	libjobs "lib/jobs"
	"lib/telemetry"
	"payssage/models"
	"payssage/ports"
)

// FailWebhookDeliveryArgs records the give-up of a webhook.deliver job
// whose retries ran out (FailDelivery), on every carrier.
type FailWebhookDeliveryArgs struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
}

func (FailWebhookDeliveryArgs) Kind() string { return "webhook.fail" }

// FailDelivery is webhook.deliver's give-up (jobs.GiveUpWith).
func FailDelivery(job DeliverWebhookArgs) libjobs.Job {
	return FailWebhookDeliveryArgs(job)
}

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
