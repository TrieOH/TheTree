package commands

import (
	"context"

	"lib/database"
	"lib/telemetry"
	"payssage/internal/features/webhooks/jobs"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// dispatchDeliveries creates a pending delivery row and enqueues a River
// job for every webhook endpoint registered on the event's wallet. If a
// wallet has no endpoints configured, this is a no-op — not an error,
// since not every wallet will have webhooks set up.
//
// The delivery row insert and the River job enqueue happen inside the
// same DB transaction (via river.InsertTx against the tx pulled from
// context), so a delivery row is never created without a job to process
// it, and vice versa — either both commit or neither does.
//
// One consequence of running the whole batch inside a single transaction:
// if any endpoint's delivery/job insert fails, the entire batch for this
// event rolls back, including endpoints that would have succeeded. This
// is a deliberate tradeoff — atomicity ("all deliveries for this event
// exist together or none do") over partial-success resilience. The
// failure path here is a DB/River insert error, not an actual delivery
// attempt failing (that's handled entirely inside the worker, with its
// own per-delivery retry budget) — so this should be a rare, infra-level
// failure mode, not a routine one.
func (c *Commands) dispatchDeliveries(ctx context.Context, event *models.WebhookEvent) error {
	ctx, span := telemetry.StartSpan(ctx, "dispatchDeliveries")
	defer span.End()

	endpoints, err := c.endpoints.ListByWallet(ctx, event.WalletID)
	if err != nil {
		return fun.Errf("list webhook endpoints: %v", err).Internal()
	}
	if len(endpoints) == 0 {
		return nil
	}

	return database.RunTx(ctx, func(ctx context.Context) error {
		tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx)
		if !ok {
			return fun.Err("dispatchDeliveries: no transaction in context").Internal()
		}

		for _, endpoint := range endpoints {
			v7, err := uuid.NewV7()
			if err != nil {
				return fun.ErrInternal(err.Error())
			}

			delivery := models.WebhookDelivery{
				ID:         v7,
				EndpointID: endpoint.ID,
				EventID:    event.ID,
				Status:     models.WebhookDeliveryStatusPending,
			}

			created, err := c.deliveries.Create(ctx, delivery)
			if err != nil {
				if fun.Is(err, fun.CodeConflict) {
					// uniq_webhook_deliveries_event_endpoint — this event
					// was already dispatched to this endpoint. Not an
					// error, just skip enqueueing a duplicate job.
					continue
				}
				return fun.Errf("create webhook delivery: %v", err).Internal()
			}

			_, err = c.river.InsertTx(ctx, tx, jobs.DeliverWebhookArgs{
				DeliveryID: created.ID,
			}, nil)
			if err != nil {
				return fun.Errf("enqueue webhook delivery job: %v", err).Internal()
			}
		}

		return nil
	})
}
