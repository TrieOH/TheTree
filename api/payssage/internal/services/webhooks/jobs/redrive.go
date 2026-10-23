package jobs

import (
	"context"
	"errors"
	"time"

	libjobs "lib/jobs"
	"lib/telemetry"
	"payssage/ports"

	"go.uber.org/zap"
)

// RedriveWebhookDeliveriesArgs re-enqueues deliveries whose enqueue was
// lost. Dispatch commits the delivery rows first and enqueues after, so a
// crash or a failed enqueue in between leaves a pending row no job will
// ever pick up; this periodic sweep finds those and enqueues them again.
type RedriveWebhookDeliveriesArgs struct{}

func (RedriveWebhookDeliveriesArgs) Kind() string { return "webhook.redrive" }

const (
	// redriveAfter is how old an unattempted delivery must be before the
	// sweep assumes its enqueue was lost; long enough that a delivery whose
	// job is only queued is not enqueued twice.
	redriveAfter = 2 * time.Minute
	// redriveBatch bounds one sweep; the next one picks up the rest.
	redriveBatch = 100
)

type RedriveWebhookDeliveriesWorker struct {
	deliveries ports.WebhookDeliveryRepo
	enqueuer   libjobs.Enqueuer
}

func NewRedriveWebhookDeliveriesWorker(deliveries ports.WebhookDeliveryRepo, enqueuer libjobs.Enqueuer) *RedriveWebhookDeliveriesWorker {
	return &RedriveWebhookDeliveriesWorker{deliveries: deliveries, enqueuer: enqueuer}
}

func (w *RedriveWebhookDeliveriesWorker) Work(ctx context.Context, _ RedriveWebhookDeliveriesArgs) error {
	unsent, err := w.deliveries.ListUnsent(ctx, redriveAfter, redriveBatch)
	if err != nil {
		return err
	}

	var errs []error
	for _, delivery := range unsent {
		err = w.enqueuer.Enqueue(ctx, DeliverWebhookArgs{DeliveryID: delivery.ID})
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(unsent) > 0 {
		telemetry.Log().Warn("webhook delivery: redrove unsent deliveries",
			zap.Int("found", len(unsent)),
			zap.Int("failed", len(errs)),
		)
	}
	return errors.Join(errs...)
}
