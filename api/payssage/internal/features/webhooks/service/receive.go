package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"lib/database"
	"lib/telemetry"
	"payssage/errs"
	"payssage/internal/providers"
	"payssage/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
