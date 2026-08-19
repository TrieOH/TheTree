package webhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"lib/telemetry"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// webhookEventTypes — the normalized event types payssage delivers
// (payment.<intent status>; see payssage's NormalizeStatus). MP's own
// vocabulary never reaches this switch. `payment.approved`/`payment.failed`
// are accepted for forward-compat with other providers.
const (
	eventSucceeded = "payment.succeeded"
	eventApproved  = "payment.approved"
	eventFailed    = "payment.failed"
	eventRejected  = "payment.rejected"
	eventCancelled = "payment.cancelled"
	eventPending   = "payment.pending"
	eventRefunded  = "payment.refunded"
)

// Receive handles one Payssage webhook delivery: verify the signature,
// correlate by intent_id (D3 card race), and dispatch on event_type.
// Returning nil acknowledges with 200 (Payssage stops retrying); returning
// an error is a non-2xx (Payssage retries).
func (o *Operations) Receive(ctx context.Context, input ReceiveInput) error {
	ctx, span := telemetry.StartSpan(ctx, "WebhooksService.Receive")
	defer span.End()

	if !o.validSignature(input) {
		telemetry.Log().Warn("webhook: invalid signature",
			zap.String("intent_id", input.IntentID.String()))
		return fun.ErrBadRequest("invalid webhook signature")
	}

	purchase, err := o.findPurchase(ctx, input.IntentID)
	if err != nil {
		return err // non-2xx → Payssage retries (up to 5 attempts ≈ 5s)
	}

	switch input.EventType {
	case eventSucceeded, eventApproved:
		return o.approve(ctx, purchase)
	case eventFailed:
		return o.cancel(ctx, purchase, models.PurchaseStatusFailed, input.StatusReason)
	case eventRejected:
		return o.cancel(ctx, purchase, models.PurchaseStatusRejected, input.StatusReason)
	case eventCancelled:
		return o.cancel(ctx, purchase, models.PurchaseStatusCancelled, input.StatusReason)
	case eventPending:
		// Non-terminal: pix's first delivery races the checkout commit (D3)
		// — by the time it resolves the purchase is pending and nothing
		// needs flipping. Acknowledge.
		return nil
	case eventRefunded:
		// Refund (refund plan slice 3, replaces the D1 deferral): flip the
		// approved purchase to refunded, cancel the materialized rows, and
		// revoke emitted badges — webhook-confirmed, like approval.
		return o.refund(ctx, purchase)
	default:
		telemetry.Log().Warn("webhook: unhandled event type",
			zap.String("event_type", input.EventType),
			zap.String("intent_id", input.IntentID.String()))
		return nil
	}
}

// validSignature verifies X-Payssage-Signature = hex(HMAC-SHA256(secret,
// raw body)) — the exact bytes payssage POSTed (D2).
func (o *Operations) validSignature(input ReceiveInput) bool {
	mac := hmac.New(sha256.New, []byte(o.secret))
	mac.Write(input.RawBody)
	return hmac.Equal([]byte(input.Signature), []byte(hex.EncodeToString(mac.Sum(nil))))
}

// findPurchase correlates the intent to a purchase (D2). Cards charge
// synchronously and pix fires its first webhook during the checkout request,
// so the delivery can beat the checkout tx commit (D3): wait ~1s and re-query
// once. If it still isn't there, return an error → non-2xx → Payssage
// retries; persistent failure is logged loudly (genuinely broken, or an
// unknown intent — e.g. a future donation intent).
func (o *Operations) findPurchase(ctx context.Context, intentID uuid.UUID) (*models.Purchase, error) {
	purchase, err := o.purchases.GetByIntentID(ctx, intentID)
	if err == nil {
		return purchase, nil
	}
	if !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}

	telemetry.Log().Info("webhook: intent not found — waiting for checkout commit (D3)",
		zap.String("intent_id", intentID.String()))
	select {
	case <-ctx.Done():
		return nil, fun.ErrInternal("webhook processing cancelled")
	case <-time.After(o.cardRaceWait):
	}

	purchase, err = o.purchases.GetByIntentID(ctx, intentID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			telemetry.Log().Error("webhook: intent never correlated to a purchase (unknown intent or broken checkout)",
				zap.String("intent_id", intentID.String()))
			return nil, fun.Errf("purchase not found for intent %s", intentID).Internal()
		}
		return nil, err
	}
	return purchase, nil
}
