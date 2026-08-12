// Package jobs holds the checkout-side River workers (split 7): the
// purchases.expire job that frees unpaid reservations.
package jobs

import (
	"context"
	"encoding/json"

	"lib/database"
	"lib/telemetry"

	"univents/internal/services/notify"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// ExpirePurchaseArgs is the `purchases.expire` job: expire one pending
// purchase at its 10:01 deadline — flip it to `expired`, flip the
// materialized rows, and notify (stock freed + purchase.expired). Does NOT
// call Payssage cancel (MP pix self-expires at 30min; cards already charged
// — late approval is the safety net).
type ExpirePurchaseArgs struct {
	PurchaseID uuid.UUID `json:"purchase_id"`
}

func (ExpirePurchaseArgs) Kind() string { return "purchases.expire" }

// Notifier is the LISTEN/NOTIFY publish surface. Satisfied by
// *database.Notifier (lib/go/database).
type Notifier interface {
	Notify(ctx context.Context, channel, payload string) error
}

// ExpirePurchaseWorker expires unpaid purchases. Idempotent: it re-checks
// the purchase status first (approve cancels the job best-effort; this
// worker is the backstop — a purchase already approved/cancelled/expired is
// a no-op). Double-run safe by the same guard.
type ExpirePurchaseWorker struct {
	river.WorkerDefaults[ExpirePurchaseArgs]

	purchases        ports.PurchaseRepo
	registrations    ports.RegistrationRepo
	productPurchases ports.ProductPurchaseRepo
	participations   ports.ProgramParticipationRepo
	notifier         Notifier
	tx               database.TxRunner
}

func NewExpirePurchaseWorker(
	purchases ports.PurchaseRepo,
	registrations ports.RegistrationRepo,
	productPurchases ports.ProductPurchaseRepo,
	participations ports.ProgramParticipationRepo,
	notifier Notifier,
	tx database.TxRunner,
) *ExpirePurchaseWorker {
	return &ExpirePurchaseWorker{
		purchases:        purchases,
		registrations:    registrations,
		productPurchases: productPurchases,
		participations:   participations,
		notifier:         notifier,
		tx:               tx,
	}
}

// Work runs the expiry: guarded pending→expired inside a tx, materialized
// rows flipped (registrations → expired, product_purchases → expired,
// program_participations → cancelled — there is no expired status), NOTIFY
// after commit. Returns nil for already-terminal purchases.
func (w *ExpirePurchaseWorker) Work(ctx context.Context, job *river.Job[ExpirePurchaseArgs]) error {
	ctx, span := telemetry.StartSpan(ctx, "ExpirePurchaseWorker.Work")
	defer span.End()

	purchase, err := w.purchases.GetByID(ctx, job.Args.PurchaseID)
	if err != nil {
		return err // missing purchase is genuinely broken — let River retry
	}
	if purchase.Status != models.PurchaseStatusPending {
		// Already approved/cancelled/expired — approve cancelled the job
		// best-effort; this run is the backstop and has nothing to do.
		telemetry.Log().Info("expiry: purchase not pending — no-op",
			zap.String("purchase_id", purchase.ID.String()),
			zap.String("status", string(purchase.Status)))
		return nil
	}

	items, err := w.purchases.ListItemsByPurchase(ctx, purchase.ID)
	if err != nil {
		return err
	}

	var stockIDs []uuid.UUID
	for _, item := range items {
		stockIDs = append(stockIDs, item.ItemID)
	}

	err = w.tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := w.purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusPending, models.PurchaseStatusExpired, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // guard missed — concurrently approved/expired; nothing to flip
		}
		return w.flipExpired(ctx, items)
	})
	if err != nil {
		return err
	}

	w.notify(ctx, purchase, stockIDs)
	return nil
}

// flipExpired materializes the expiry (D4), inside the expire tx:
// registrations → expired, product_purchases → expired,
// program_participations → cancelled (no expired status exists for them).
func (w *ExpirePurchaseWorker) flipExpired(ctx context.Context, items []models.PurchaseItem) error {
	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				continue
			}
			_, err := w.registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusExpired, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				continue
			}
			_, err := w.productPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusExpired, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				continue
			}
			_, err := w.participations.UpdateParticipationStatus(ctx, *item.ParticipationID, models.ProgramParticipationStatusCancelled)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// notify publishes the freed stock deltas (D10) and the purchase.expired
// event (D9). Fire-and-forget: a missed notification is a stale snapshot,
// never data loss.
func (w *ExpirePurchaseWorker) notify(ctx context.Context, purchase *models.Purchase, stockIDs []uuid.UUID) {
	stock := notify.StockNotification{Kind: notify.KindStock, EditionID: purchase.EditionID, ItemIDs: stockIDs}
	w.publish(ctx, stock)
	w.publish(ctx, notify.PurchaseNotification{
		Kind:       notify.KindPurchase,
		EditionID:  purchase.EditionID,
		PurchaseID: purchase.ID,
		Status:     string(models.PurchaseStatusExpired),
	})
}

func (w *ExpirePurchaseWorker) publish(ctx context.Context, payload any) {
	raw, err := json.Marshal(payload)
	if err != nil {
		telemetry.Log().Error("expiry: marshal notifier payload",
			zap.String("channel", notify.Channel),
			zap.Error(err))
		return
	}
	err = w.notifier.Notify(ctx, notify.Channel, string(raw))
	if err != nil {
		telemetry.Log().Error("expiry: publish to notifier",
			zap.String("channel", notify.Channel),
			zap.String("payload", string(raw)),
			zap.Error(err))
	}
}
