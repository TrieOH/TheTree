package service

import (
	"context"

	"lib/telemetry"
	"payssage/internal/features/webhooks/jobs"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// dispatchDeliveries creates a pending delivery row for every webhook
// endpoint registered on the event's wallet, then enqueues a delivery job
// for each. If a wallet has no endpoints configured, this is a no-op — not
// an error, since not every wallet will have webhooks set up.
//
// The rows are the source of truth; the job only tells a worker to go
// deliver one. So the rows commit together in one transaction ("all
// deliveries for this event exist together or none do"), and the jobs are
// enqueued after the commit, outside it. An enqueue that fails is logged,
// not returned: the row is safe, and the webhook.redrive sweep enqueues
// pending rows that were never attempted. Enqueueing before the commit
// would let a worker look for a row that is not visible yet.
func (o *Operations) dispatchDeliveries(ctx context.Context, event *models.WebhookEvent) error {
	ctx, span := telemetry.StartSpan(ctx, "dispatchDeliveries")
	defer span.End()

	endpoints, err := o.Endpoints.ListByWallet(ctx, event.WalletID)
	if err != nil {
		return fun.Errf("list webhook endpoints: %v", err).Internal()
	}
	if len(endpoints) == 0 {
		return nil
	}

	created := make([]uuid.UUID, 0, len(endpoints))
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		for _, endpoint := range endpoints {
			v7, err := uuid.NewV7()
			if err != nil {
				return fun.ErrInternal(err.Error())
			}

			delivery, err := o.Deliveries.Create(ctx, models.WebhookDelivery{
				ID:         v7,
				EndpointID: endpoint.ID,
				EventID:    event.ID,
				Status:     models.WebhookDeliveryStatusPending,
			})
			if err != nil {
				if fun.Is(err, fun.CodeConflict) {
					// uniq_webhook_deliveries_event_endpoint — this event
					// was already dispatched to this endpoint. Not an
					// error, just skip enqueueing a duplicate job.
					continue
				}
				return fun.Errf("create webhook delivery: %v", err).Internal()
			}
			created = append(created, delivery.ID)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, id := range created {
		err = o.Jobs.Enqueue(ctx, jobs.DeliverWebhookArgs{DeliveryID: id})
		if err != nil {
			telemetry.Log().Error("webhook delivery: enqueue failed, left for redrive",
				zap.String("delivery_id", id.String()),
				zap.Error(err),
			)
		}
	}
	return nil
}
