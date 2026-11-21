// Package webhooks is the webhooks feature: its handlers and operations.
// Build wires it.
//
// The operations are the Payssage webhook receiver — the only component
// that confirms payment (D3). It verifies the delivery signature, correlates
// the intent to a purchase (D2), and flips the purchase + materialized rows
// on approval, failure, or late approval (D4). Checkout never self-approves.
package webhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"lib/database"
	libjobs "lib/jobs"
	"lib/telemetry"
	"log/slog"
	"time"
	"univents/errs"
	"univents/internal/features/checkouts/jobs"
	"univents/internal/notify"
	"univents/models"
	"univents/ports"
	"uuid"
)

// Badges is the badge-emission surface the receiver calls on approval and
// refund. Satisfied by *badges.Operations.
type Badges interface {
	EmitForConfirmedRegistration(ctx context.Context, registrationID uuid.UUID) (*models.BadgeEmission, error)
	// RevokeForRegistration revokes a participant badge when its
	// registration is cancelled (refund plan slice 3: refunds flip
	// registrations to cancelled and revoke the emitted badge).
	RevokeForRegistration(ctx context.Context, registrationID uuid.UUID, reason string) error
}

// Notifier is the LISTEN/NOTIFY publish surface (split 6 subscribes).
// Satisfied by *database.Notifier (lib/go/database).
type Notifier interface {
	Notify(ctx context.Context, channel, payload string) error
}

// Jobs is the background-work surface: the gifted-ticket email
// (gifts.send_email) is enqueued inside the approve tx, before commit, with
// a delay so it cannot run before the confirmation is visible. Approval
// does not cancel the purchase's expiry job — a queued message cannot be
// withdrawn; the expiry worker no-ops on a purchase that is no longer
// pending. Satisfied by the lib/jobs carriers.
type Jobs interface {
	EnqueueAfter(ctx context.Context, job libjobs.Job, delay time.Duration) error
}

// channelUniventsChanges is the single NOTIFY channel the store publishes
// on (contract in internal/notify): the WS hub routes
// kind="purchase" events (D9) and the SSE relay routes kind="stock" deltas
// (D10), both subscribed in split 6.
const channelUniventsChanges = notify.Channel

type Deps struct {
	Purchases        ports.PurchaseRepo
	Registrations    ports.RegistrationRepo
	ProductPurchases ports.ProductPurchaseRepo
	Participations   ports.ProgramParticipationRepo
	Badges           Badges
	Notifier         Notifier
	Jobs             Jobs
	Tx               database.TxRunner
	Secret           string
}

type Operations struct {
	Deps

	// cardRaceWait is the D3 wait before re-querying a purchase whose
	// checkout tx may not have committed yet. Overridable in tests.
	cardRaceWait time.Duration
}

func New(d Deps) *Operations {
	return &Operations{
		Deps:         d,
		cardRaceWait: time.Second,
	}
}

// SetCardRaceWait overrides the D3 wait between the first miss and the
// re-query. Test-only (zero makes the race tests instant).
func (o *Operations) SetCardRaceWait(d time.Duration) { o.cardRaceWait = d }

// ReceiveInput is one Payssage webhook delivery, reduced to what the
// receiver acts on: correlation key (intent_id), event type, and the raw
// body + signature header for verification.
type ReceiveInput struct {
	IntentID     uuid.UUID
	EventType    string
	RawBody      []byte
	Signature    string
	StatusReason *string // normalized outcome detail from the payssage envelope (e.g. "high_risk")
}

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
		telemetry.Log().WarnContext(ctx, "webhook: approval for terminal purchase ignored",
			slog.String("purchase_id", purchase.ID.String()),
			slog.String("status", string(purchase.Status)))
		return nil
	}
}

// approvePending is the normal approve path, in a tx: guarded pending→
// approved flip (idempotency — a concurrent duplicate delivery misses the
// guard), expiry-job cancellation (best-effort), materialized rows flipped
// to confirmed (+ badge emit), and NOTIFY (stock deltas + purchase event).
func (o *Operations) approvePending(ctx context.Context, purchase *models.Purchase) error {
	return o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.Purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusPending, models.PurchaseStatusApproved, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — already flipped by a concurrent delivery
		}
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
	items, err := o.Purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}
	availability, err := o.Purchases.Availability(ctx, purchase.EditionID)
	if err != nil {
		return err
	}

	if !fullyAvailable(items, availability) {
		reason := statusReasonRefundedAfterExpiry
		updated, err := o.Purchases.UpdateStatus(ctx, purchase.ID, models.PurchaseStatusExpired, &reason)
		if err != nil {
			return err
		}
		o.notify(ctx, updated, items, models.PurchaseStatusExpired)
		telemetry.Log().ErrorContext(ctx, "webhook: late approval without full stock — refund deferred (D1), manual reconciliation needed",
			slog.String("purchase_id", purchase.ID.String()),
			slog.String("intent_id", intentIDString(purchase)),
			slog.String("edition_id", purchase.EditionID.String()))
		return nil
	}

	// Honor the order: approve from expired, same as the pending path.
	return o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.Purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusExpired, models.PurchaseStatusApproved, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — concurrently flipped (re-expired or approved)
		}
		return o.flipConfirmed(ctx, updated)
	})
}

// flipConfirmed materializes the approval (D4), inside the approve tx:
// registrations pending→confirmed (+ badge emit — no-op for non-confirmed),
// product_purchases pending→confirmed, program_participations stay
// registered (attendance is tracked against the ticket's registration).
// Then NOTIFYs stock deltas and the purchase event.
func (o *Operations) flipConfirmed(ctx context.Context, purchase *models.Purchase) error {
	items, err := o.Purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				telemetry.Log().WarnContext(ctx, "webhook: ticket item without registration",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			reg, err := o.Registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusConfirmed, nil)
			if err != nil {
				return err
			}
			_, err = o.Badges.EmitForConfirmedRegistration(ctx, *item.RegistrationID)
			if err != nil {
				return err
			}
			o.enqueueGiftEmail(ctx, reg)
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				telemetry.Log().WarnContext(ctx, "webhook: product item without product_purchase",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.ProductPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusConfirmed, nil)
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

// giftEmailDelay holds the gift email back long enough for the approve tx
// to commit: it is enqueued inside the tx, before commit.
const giftEmailDelay = 30 * time.Second

// enqueueGiftEmail schedules the gifted-ticket email for an accountless
// recipient (attendee_user_id NULL) inside the approve tx — the same tx
// that just confirmed the registration. The delay keeps the job from
// running before the confirmation commits; if the tx rolls back, the job
// finds the registration unconfirmed and no-ops. Account holders are
// skipped: they see the ticket in my-ticket and get the badge email.
func (o *Operations) enqueueGiftEmail(ctx context.Context, reg *models.Registration) {
	if reg.AttendeeUserID != nil {
		return
	}
	err := o.Jobs.EnqueueAfter(ctx, jobs.SendGiftEmailArgs{RegistrationID: reg.ID}, giftEmailDelay)
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "webhook: failed to enqueue gift email",
			slog.String("registration_id", reg.ID.String()),
			telemetry.Err(err))
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
	stock := stockNotification(purchase, items)
	o.publish(ctx, stock)
	o.publish(ctx, notify.PurchaseNotification{
		Kind:       notify.KindPurchase,
		EditionID:  purchase.EditionID,
		PurchaseID: purchase.ID,
		Status:     string(status),
	})
}

// stockNotification builds the stock-delta payload (D10): item_ids only —
// the SSE relay re-queries availability from the DB.
func stockNotification(purchase *models.Purchase, items []models.PurchaseItem) notify.StockNotification {
	stock := notify.StockNotification{Kind: notify.KindStock, EditionID: purchase.EditionID}
	for _, item := range items {
		stock.ItemIDs = append(stock.ItemIDs, item.ItemID)
	}
	return stock
}

func (o *Operations) publish(ctx context.Context, payload any) {
	raw, err := json.Marshal(payload)
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "webhook: marshal notifier payload",
			slog.String("channel", channelUniventsChanges),
			telemetry.Err(err))
		return
	}
	err = o.Notifier.Notify(ctx, channelUniventsChanges, string(raw))
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "webhook: publish to notifier",
			slog.String("channel", channelUniventsChanges),
			slog.String("payload", string(raw)),
			telemetry.Err(err))
	}
}

func intentIDString(purchase *models.Purchase) string {
	if purchase.PayssageIntentID == nil {
		return ""
	}
	return purchase.PayssageIntentID.String()
}

// cancel handles a terminal failure webhook (failed/rejected/cancelled): in
// a tx, guarded pending→target (a stale failure delivery can't cancel an
// approved purchase), materialized rows flipped to cancelled, NOTIFY (stock
// freed + purchase event). Non-pending purchases are left alone and
// acknowledged.
//
// The target status distinguishes why the purchase died (D4): the customer
// cancelling is "cancelled", a declined payment is "failed" (payment.failed)
// or "rejected" (payment.rejected — e.g. MP's risk engine). The normalized
// reason from the delivery envelope lands in status_reason.
func (o *Operations) cancel(ctx context.Context, purchase *models.Purchase, to models.PurchaseStatus, reason *string) error {
	if purchase.Status != models.PurchaseStatusPending {
		telemetry.Log().InfoContext(ctx, "webhook: failure for non-pending purchase ignored",
			slog.String("purchase_id", purchase.ID.String()),
			slog.String("status", string(purchase.Status)))
		return nil
	}
	return o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.Purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusPending, to, reason)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — already flipped by a concurrent delivery
		}
		return o.flipCancelled(ctx, updated, to)
	})
}

// flipCancelled materializes the cancellation (D4), inside the same tx:
// registrations/product_purchases/participations flipped to cancelled, then
// NOTIFY with the purchase's target status (stock freed + purchase event).
func (o *Operations) flipCancelled(ctx context.Context, purchase *models.Purchase, to models.PurchaseStatus) error {
	items, err := o.Purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				telemetry.Log().WarnContext(ctx, "webhook: ticket item without registration",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.Registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				telemetry.Log().WarnContext(ctx, "webhook: product item without product_purchase",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.ProductPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				telemetry.Log().WarnContext(ctx, "webhook: program item without participation",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.Participations.UpdateParticipationStatus(ctx, *item.ParticipationID, models.ProgramParticipationStatusCancelled)
			if err != nil {
				return err
			}
		}
	}

	o.notify(ctx, purchase, items, to)
	return nil
}

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
	if !o.validSignature(input) {
		telemetry.Log().WarnContext(ctx, "webhook: invalid signature",
			slog.String("intent_id", input.IntentID.String()))
		return errs.ErrInvalidPayssageWebhookSignature
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
		telemetry.Log().WarnContext(ctx, "webhook: unhandled event type",
			slog.String("event_type", input.EventType),
			slog.String("intent_id", input.IntentID.String()))
		return nil
	}
}

// validSignature verifies X-Payssage-Signature = hex(HMAC-SHA256(secret,
// raw body)) — the exact bytes payssage POSTed (D2).
func (o *Operations) validSignature(input ReceiveInput) bool {
	mac := hmac.New(sha256.New, []byte(o.Secret))
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
	purchase, err := o.Purchases.GetByIntentID(ctx, intentID)
	if err == nil {
		return purchase, nil
	}
	if !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}

	telemetry.Log().InfoContext(ctx, "webhook: intent not found — waiting for checkout commit (D3)",
		slog.String("intent_id", intentID.String()))
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("webhook processing cancelled: %w", ctx.Err())
	case <-time.After(o.cardRaceWait):
	}

	purchase, err = o.Purchases.GetByIntentID(ctx, intentID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			telemetry.Log().ErrorContext(ctx, "webhook: intent never correlated to a purchase (unknown intent or broken checkout)",
				slog.String("intent_id", intentID.String()))
			return nil, fmt.Errorf("purchase not found for intent %s", intentID)
		}
		return nil, err
	}
	return purchase, nil
}

// refund handles a payment.refunded webhook (refund plan slice 3): in a tx,
// guarded approved→refunded (a duplicate refund delivery no-ops), materialized
// rows flipped to cancelled, emitted badges revoked, and stock NOTIFYed (the
// availability ledger only counts pending/approved, so leaving approved frees
// the items automatically). Non-approved purchases are left alone and
// acknowledged.
func (o *Operations) refund(ctx context.Context, purchase *models.Purchase) error {
	if purchase.Status != models.PurchaseStatusApproved {
		telemetry.Log().InfoContext(ctx, "webhook: refund for non-approved purchase ignored",
			slog.String("purchase_id", purchase.ID.String()),
			slog.String("status", string(purchase.Status)))
		return nil
	}
	return o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.Purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusApproved, models.PurchaseStatusRefunded, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — already flipped by a concurrent delivery
		}
		return o.flipRefunded(ctx, updated)
	})
}

// flipRefunded materializes the refund (slice 3), inside the same tx:
// registrations/product_purchases/participations flipped to cancelled and the
// ticket badges revoked (a refunded buyer must not keep a badge), then NOTIFY
// stock deltas (item ids) so the SSE relay re-reads availability. The D9
// purchase event is deliberately NOT published — the buyer's WS socket closed
// on approval and the hub's framesFor default logs unknown statuses; there is
// no subscriber to a refunded purchase frame.
func (o *Operations) flipRefunded(ctx context.Context, purchase *models.Purchase) error {
	items, err := o.Purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				telemetry.Log().WarnContext(ctx, "webhook: refunded ticket item without registration",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.Registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusCancelled, nil)
			if err != nil {
				return err
			}
			err = o.Badges.RevokeForRegistration(ctx, *item.RegistrationID, "refunded")
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				telemetry.Log().WarnContext(ctx, "webhook: refunded product item without product_purchase",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.ProductPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				telemetry.Log().WarnContext(ctx, "webhook: refunded program item without participation",
					slog.String("purchase_id", purchase.ID.String()))
				continue
			}
			_, err := o.Participations.UpdateParticipationStatus(ctx, *item.ParticipationID, models.ProgramParticipationStatusCancelled)
			if err != nil {
				return err
			}
		}
	}

	o.notifyStockOnly(ctx, purchase, items)
	return nil
}

// notifyStockOnly publishes the refunded purchase's stock deltas (D10) —
// item ids only, the SSE relay recomputes from the DB. Unlike approve/cancel
// there is no purchase event: no subscriber exists for a refunded frame.
func (o *Operations) notifyStockOnly(ctx context.Context, purchase *models.Purchase, items []models.PurchaseItem) {
	stock := stockNotification(purchase, items)
	o.publish(ctx, stock)
}
