package webhooks

import (
	"context"
	"lib/telemetry"
	"payssage/internal/providers"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (o *Operations) Receive(ctx context.Context, payload models.ReceiveWebhookInput) error {
	ctx, span := telemetry.StartSpan(ctx, "ReceiveWebhook")
	defer span.End()

	providerEnum, err := providers.FromString(payload.Provider)
	if err != nil {
		return fun.Err("unknown webhook provider").NotFound()
	}
	handler, ok := providers.PayssageProviders.Webhooks[providerEnum]
	if !ok {
		return fun.Err("webhook handling not implemented for this provider").BadRequest()
	}

	err = handler.VerifySignature(ctx, payload.Request, payload.RawBody)
	if err != nil {
		return err // stays a real error — bad signature should NOT return 200
	}

	result, err := handler.Parse(ctx, payload.Request, payload.RawBody)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			// Covers both "no matching intent" and "unhandled event type" —
			// both are expected-and-ignored cases, not failures. Acknowledge
			// with 200 so the provider doesn't retry.
			telemetry.Log().Info("webhook acknowledged, not processed",
				zap.String("provider", payload.Provider),
				zap.Error(err),
			)
			return nil
		}
		return err
	}
	event := models.WebhookEvent{
		ID:           uuid.Must(uuid.NewV7()),
		WalletID:     result.WalletID,
		IntentID:     result.IntentID,
		Provider:     payload.Provider,
		ExternalID:   result.ExternalID,
		EventType:    result.EventType,
		StatusDetail: result.StatusDetail,
		Payload:      payload.RawBody,
	}

	created, err := o.events.Create(ctx, event)
	if err != nil {
		if fun.Is(err, fun.CodeConflict) {
			// Same (provider, external_id, event_type) already recorded and
			// dispatched — an identical redelivery, nothing further to do.
			// Note: DISTINCT event types for the same payment are NOT
			// duplicates (payment.pending vs payment.succeeded) — the dedupe
			// index keys on event_type too (migration 007), so the final
			// status always dispatches.
			return nil
		}
		return fun.Errf("persist webhook event: %v", err).Internal()
	}

	return o.dispatchDeliveries(ctx, created)
}
