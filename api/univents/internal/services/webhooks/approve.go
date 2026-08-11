package webhooks

import (
	"context"
	"encoding/json"

	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// statusReasonRefundedAfterExpiry marks a late approval that could not be
// honored (D1): the purchase stays expired and the money is refunded
// manually, off-band, until the refund endpoint ships.
const statusReasonRefundedAfterExpiry = "refunded_after_expiry"

// approve confirms a purchase paid via the webhook — the only payment
// confirmation path (D3). Idempotent: a duplicate delivery finds the
// purchase already approved and becomes a no-op.
func (o *Operations) approve(ctx context.Context, purchase *models.Purchase) error {
	switch purchase.Status {
	case models.PurchaseStatusApproved:
		return nil // duplicate delivery — already approved
	case models.PurchaseStatusPending:
		return o.approvePending(ctx, purchase)
	case models.PurchaseStatusExpired:
		return o.lateApprove(ctx, purchase)
	default: // cancelled — terminal
		telemetry.Log().Warn("webhook: approval for terminal purchase ignored",
			zap.String("purchase_id", purchase.ID.String()),
			zap.String("status", string(purchase.Status)))
		return nil
	}
}

// approvePending is the normal approve path, in a tx: guarded pending→
// approved flip (idempotency — a concurrent duplicate delivery misses the
// guard), expiry-job cancellation (best-effort), materialized rows flipped
// to confirmed (+ badge emit), and NOTIFY (stock deltas + purchase event).
func (o *Operations) approvePending(ctx context.Context, purchase *models.Purchase) error {
	return o.tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusPending, models.PurchaseStatusApproved, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — already flipped by a concurrent delivery
		}
		o.cancelExpiryJob(ctx, updated)
		return o.flipConfirmed(ctx, updated)
	})
}

// lateApprove handles an approval that arrives after the 10:01 expiry (the
// buyer's socket already closed at expiry; a reconnect shows the resolved
// state). Re-check availability of every item (split 3 semantics): fully
// available → honor the order (approve as above); otherwise → refund
// deferred (D1): keep expired, set status_reason, NOTIFY, and log loudly
// for manual reconciliation.
func (o *Operations) lateApprove(ctx context.Context, purchase *models.Purchase) error {
	items, err := o.purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}
	availability, err := o.purchases.Availability(ctx, purchase.EditionID)
	if err != nil {
		return err
	}

	if !fullyAvailable(items, availability) {
		reason := statusReasonRefundedAfterExpiry
		updated, err := o.purchases.UpdateStatus(ctx, purchase.ID, models.PurchaseStatusExpired, &reason)
		if err != nil {
			return err
		}
		o.notify(ctx, updated, items, models.PurchaseStatusExpired)
		telemetry.Log().Error("webhook: late approval without full stock — refund deferred (D1), manual reconciliation needed",
			zap.String("purchase_id", purchase.ID.String()),
			zap.String("intent_id", intentIDString(purchase)),
			zap.String("edition_id", purchase.EditionID.String()))
		return nil
	}

	// Honor the order: approve from expired, same as the pending path.
	return o.tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusExpired, models.PurchaseStatusApproved, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — concurrently flipped (re-expired or approved)
		}
		o.cancelExpiryJob(ctx, updated)
		return o.flipConfirmed(ctx, updated)
	})
}

// flipConfirmed materializes the approval (D4), inside the approve tx:
// registrations pending→confirmed (+ badge emit — no-op for non-confirmed),
// product_purchases pending→confirmed, program_participations stay
// registered (attendance is tracked against the ticket's registration).
// Then NOTIFYs stock deltas and the purchase event.
func (o *Operations) flipConfirmed(ctx context.Context, purchase *models.Purchase) error {
	items, err := o.purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				telemetry.Log().Warn("webhook: ticket item without registration",
					zap.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusConfirmed, nil)
			if err != nil {
				return err
			}
			_, err = o.badges.EmitForConfirmedRegistration(ctx, *item.RegistrationID)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				telemetry.Log().Warn("webhook: product item without product_purchase",
					zap.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.productPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusConfirmed, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			// program_participations stay registered on approval (D4).
		}
	}

	o.notify(ctx, purchase, items, models.PurchaseStatusApproved)
	return nil
}

// cancelExpiryJob cancels the 10:01 expiry job when the purchase has one
// (split 7 checkout schedules it; seeded split-4 purchases don't). Best-
// effort: the expiry worker re-checks status before expiring, so a missed
// cancellation is harmless.
func (o *Operations) cancelExpiryJob(ctx context.Context, purchase *models.Purchase) {
	if purchase.RiverJobID == nil {
		return
	}
	_, err := o.river.JobCancel(ctx, *purchase.RiverJobID)
	if err != nil {
		telemetry.Log().Warn("webhook: failed to cancel expiry job",
			zap.String("purchase_id", purchase.ID.String()),
			zap.Int64("river_job_id", *purchase.RiverJobID),
			zap.Error(err))
	}
}

// fullyAvailable reports whether every item in the purchase still has stock
// for its quantity (available = base - reserved; nil base = unlimited). The
// expired purchase itself is excluded from reserved — availability only
// counts pending/approved purchases — so this sees other active claims.
func fullyAvailable(items []models.PurchaseItem, availability []models.ItemAvailability) bool {
	stock := make(map[uuid.UUID]models.ItemAvailability, len(availability))
	for _, a := range availability {
		stock[a.ItemID] = a
	}
	for _, item := range items {
		a, ok := stock[item.ItemID]
		if !ok {
			return false // item missing from the availability read — treat as unavailable
		}
		if a.BaseQuantity == nil {
			continue // unlimited — never sold out
		}
		if int64(*a.BaseQuantity)-a.ReservedQuantity < int64(item.Quantity) {
			return false
		}
	}
	return true
}

// notify publishes the purchase's stock deltas (D10) and its status event
// (D9) on univents_changes. Fire-and-forget: the SSE relay re-queries the
// DB and the WS hub routes on purchase_id, so a missed notification is a
// stale snapshot, never data loss — errors are logged, not propagated.
func (o *Operations) notify(ctx context.Context, purchase *models.Purchase, items []models.PurchaseItem, status models.PurchaseStatus) {
	stock := stockNotification{Kind: kindStock, EditionID: purchase.EditionID}
	for _, item := range items {
		stock.ItemIDs = append(stock.ItemIDs, item.ItemID)
	}
	o.publish(ctx, stock)
	o.publish(ctx, purchaseNotification{
		Kind:       kindPurchase,
		EditionID:  purchase.EditionID,
		PurchaseID: purchase.ID,
		Status:     string(status),
	})
}

func (o *Operations) publish(ctx context.Context, payload any) {
	raw, err := json.Marshal(payload)
	if err != nil {
		telemetry.Log().Error("webhook: marshal notifier payload",
			zap.String("channel", channelUniventsChanges),
			zap.Error(err))
		return
	}
	err = o.notifier.Notify(ctx, channelUniventsChanges, string(raw))
	if err != nil {
		telemetry.Log().Error("webhook: publish to notifier",
			zap.String("channel", channelUniventsChanges),
			zap.String("payload", string(raw)),
			zap.Error(err))
	}
}

func intentIDString(purchase *models.Purchase) string {
	if purchase.PayssageIntentID == nil {
		return ""
	}
	return purchase.PayssageIntentID.String()
}
