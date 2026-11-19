// Package webhooks is the webhooks feature: its handlers and operations.
// Build wires it.
package webhooks

import (
	"context"
	"errors"
	"fmt"
	"lib/database"
	libjobs "lib/jobs"
	"lib/telemetry"
	"log/slog"
	"net/http"
	"payssage/errs"
	"payssage/internal/authz"
	"payssage/internal/features/webhooks/jobs"
	"payssage/internal/providers"
	"payssage/models"
	"payssage/ports"
	"uuid"
)

type Deps struct {
	Jobs       libjobs.Enqueuer
	Events     ports.WebhookEventRepo
	Endpoints  ports.WebhookEndpointRepo
	Deliveries ports.WebhookDeliveryRepo
	Authz      *authz.Service
	Registry   providers.PaymentProviders
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

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
	endpoints, err := o.Endpoints.ListByWallet(ctx, event.WalletID)
	if err != nil {
		return fmt.Errorf("list webhook endpoints: %w", err)
	}
	if len(endpoints) == 0 {
		return nil
	}

	created := make([]uuid.UUID, 0, len(endpoints))
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		for _, endpoint := range endpoints {
			v7 := uuid.NewV7()

			delivery, err := o.Deliveries.Create(ctx, models.WebhookDelivery{
				ID:         v7,
				EndpointID: endpoint.ID,
				EventID:    event.ID,
				Status:     models.WebhookDeliveryStatusPending,
			})
			if err != nil {
				if errors.Is(err, database.ErrUniqueViolation) {
					// uniq_webhook_deliveries_event_endpoint — this event
					// was already dispatched to this endpoint. Not an
					// error, just skip enqueueing a duplicate job.
					continue
				}
				return fmt.Errorf("create webhook delivery: %w", err)
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
				slog.String("delivery_id", id.String()),
				telemetry.Err(err),
			)
		}
	}
	return nil
}

func (o *Operations) Receive(ctx context.Context, payload models.ReceiveWebhookInput) error {
	providerEnum, err := providers.FromString(payload.Provider)
	if err != nil {
		return errs.ErrUnknownWebhookProvider
	}
	handler, ok := o.Registry.Webhooks[providerEnum]
	if !ok {
		return errs.ErrWebhookNotImplemented
	}

	err = handler.VerifySignature(ctx, payload.Request, payload.RawBody)
	if err != nil {
		return err // stays a real error — bad signature should NOT return 200
	}

	result, err := handler.Parse(ctx, payload.Request, payload.RawBody)
	if err != nil {
		if ignorable(err) {
			// Covers "no matching intent", "unhandled event type" and a
			// payment the provider no longer has — expected-and-ignored
			// cases, not failures. Acknowledge with 200 so the provider
			// doesn't retry.
			telemetry.Log().Info("webhook acknowledged, not processed",
				slog.String("provider", payload.Provider),
				telemetry.Err(err),
			)
			return nil
		}
		return err
	}
	event := models.WebhookEvent{
		ID:           uuid.NewV7(),
		WalletID:     result.WalletID,
		IntentID:     result.IntentID,
		Provider:     payload.Provider,
		ExternalID:   result.ExternalID,
		EventType:    result.EventType,
		StatusDetail: result.StatusDetail,
		Payload:      payload.RawBody,
	}

	created, err := o.Events.Create(ctx, event)
	if err != nil {
		if errors.Is(err, database.ErrUniqueViolation) {
			// Same (provider, external_id, event_type) already recorded and
			// dispatched — an identical redelivery, nothing further to do.
			// Note: DISTINCT event types for the same payment are NOT
			// duplicates (payment.pending vs payment.succeeded) — the dedupe
			// index keys on event_type too (migration 007), so the final
			// status always dispatches.
			return nil
		}
		return fmt.Errorf("persist webhook event: %w", err)
	}

	return o.dispatchDeliveries(ctx, created)
}

// ignorable reports whether a webhook Parse failure is one to acknowledge
// and drop rather than fail (which makes the provider retry).
func ignorable(err error) bool {
	if errors.Is(err, database.ErrNotFound) {
		return true
	}
	var unhandled *errs.UnhandledWebhookEventError
	var noIntent *errs.WebhookIntentNotFoundError
	var perr *errs.ProviderError
	return errors.As(err, &unhandled) ||
		errors.As(err, &noIntent) ||
		(errors.As(err, &perr) && perr.Status == http.StatusNotFound)
}
