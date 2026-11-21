package checkouts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lib/database"
	"lib/fielderr"
	"lib/telemetry"
	"log/slog"
	"time"
	"univents/errs"
	"univents/internal/features/checkouts/jobs"
	"univents/internal/notify"
	"univents/models"
	"uuid"
)

// reservationDuration is how long an unpaid reservation is held: 10:01 (the
// master plan's expiry budget). The purchase's expires_at and the expiry
// job's delay share this one value.
const reservationDuration = 10*time.Minute + time.Second

// Payment method vocabulary from the checkout request (mirrors the spec).
const (
	paymentPix        = "pix"
	paymentCreditCard = "credit_card"
	paymentDebitCard  = "debit_card"
	currencyBRL       = "BRL"
)

// ── Domain input ──────────────────────────────────────────────────────────

// CheckoutInput is the domain form of the createEditionCheckout request.
// Prices are never trusted from the client — the server recomputes them
// from the DB inside the transaction.
type CheckoutInput struct {
	PaymentMethod   string
	CardToken       *string
	PaymentMethodID *string
	IssuerID        *string // MP issuing bank id from the front's card tokenization; optional, cards only
	Installments    *int
	Payer           Payer
	Items           []CheckoutLine
}

// Payer is the payer identity passed through to the provider (MercadoPago
// requires email + CPF/CNPJ identification).
type Payer struct {
	Email                string
	IdentificationType   string
	IdentificationNumber string
}

// CheckoutLine is one requested cart line. Ticket lines are one-per-person
// (quantity 1, one attendee); product lines carry quantity; program lines
// are quantity 1 (one session) and attach to a ticket's registration.
type CheckoutLine struct {
	ItemType models.PurchaseItemType
	ItemID   uuid.UUID
	Quantity int
	Attendee *Attendee // required for tickets
}

// Attendee is the person a ticket unit is assigned to — the purchaser
// themselves or a gifted recipient. Email is always required; UserID is
// optional: the server resolves the email against IdentityX at checkout,
// ties the account's actor id when one exists, and leaves it nil for
// email-only gifts (the recipient has no account yet). When both are sent,
// the id must belong to the email (verified server-side).
type Attendee struct {
	UserID *uuid.UUID
	Email  string
	Name   string
}

// CheckoutResult is what a successful checkout hands back: the created
// purchase (mirroring the shared Purchase read shape) plus the one-time WS
// handshake token so the front can open the socket immediately.
type CheckoutResult struct {
	PurchaseID       uuid.UUID
	EditionID        uuid.UUID
	Status           models.PurchaseStatus
	StatusReason     *string
	TotalCents       int64
	Currency         string
	PaymentMethod    *string
	PayssageIntentID *uuid.UUID
	QRCode           *string
	QRCodeBase64     *string
	ExpiresAt        time.Time
	CreatedAt        time.Time
	Items            []models.PurchaseItem
	WsToken          string
	WsTokenExpiresAt time.Time
}

// ── Checkout ──────────────────────────────────────────────────────────────

// Checkout is the money path (split 7): reserve the cart and create the
// Payssage intent in one request. Synchronous; a single transaction
// materializes the pending purchase (+ expiry job + ws token) and commits,
// then the intent is created against the event's seller on the platform
// wallet, and the intent id is stored back in a second transaction. Only
// the webhook receiver confirms payment (D3) — checkout never self-approves
// (cards charge synchronously; the purchase stays pending until the webhook
// confirms). Free orders (total 0) confirm immediately without an intent.
func (o *Operations) Checkout(ctx context.Context, editionID, purchaserID uuid.UUID, in CheckoutInput) (*CheckoutResult, error) {
	// 1. Edition + event + payment config (plain reads — the tx only locks
	//    the item rows).
	edition, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}
	if edition.IsDraft {
		return nil, errs.ErrEditionNotPublished
	}
	event, err := o.Events.GetByID(ctx, edition.EventID)
	if err != nil {
		return nil, err
	}

	// 2. Validate the cart (shape errors → 400 listing the offending lines)
	//    and the payment payload (cards need token + payment_method_id).
	lines, err := o.validateLines(in.Items)
	if err != nil {
		return nil, err
	}

	// 2b. Resolve every ticket attendee's identity against IdentityX
	//    (HTTP — outside the tx): tie the account's actor id when the email
	//    has one, verify id+email pairs, dedupe the cart. Mismatches and
	//    unknown emails with a claimed id → 400.
	lines, err = o.resolveAttendees(ctx, lines)
	if err != nil {
		return nil, err
	}

	// 3. Tx 1: lock the item rows, recompute prices, check availability,
	//    create the purchase + items + materialized rows, schedule the
	//    expiry job (delayed, enqueued before commit), issue the ws token — commit.
	res, err := o.reserve(ctx, editionID, purchaserID, lines, in)
	if err != nil {
		return nil, err
	}

	// The reservation is visible — tell the SSE relay to re-read stock
	// (D10: item_ids only, the relay recomputes from the DB).
	o.notifyStock(ctx, res.purchase, res.items)

	if res.free {
		// Free order: already approved + materialized confirmed in the tx.
		// The hub turns the purchase event into the confirmed frame.
		o.notifyPurchase(ctx, res.purchase, models.PurchaseStatusApproved)
		return o.result(res), nil
	}
	if event.PayssageSellerID == nil || event.PayssagePublicKey == nil {
		o.revert(ctx, res.purchase, res.items)
		return nil, errs.ErrNoPaymentConfig
	}
	err = o.validatePayment(in)
	if err != nil {
		o.revert(ctx, res.purchase, res.items)
		return nil, err
	}

	// 4. Create the Payssage intent (post-commit — no HTTP inside the DB
	//    tx; if this fails the purchase is reverted so the buyer can retry
	//    immediately).
	intent, err := o.Payssage.Checkout(ctx, o.WalletID, o.buildIntentRequest(res.purchase, in, *event.PayssageSellerID, res.priced))
	if err != nil {
		o.revert(ctx, res.purchase, res.items)
		return nil, mapPayssageError(err)
	}

	// 5. Tx 2: store intent id + seller + QR on the purchase.
	err = o.attachIntent(ctx, res, intent)
	if err != nil {
		// Rare infra failure: best-effort cancel the intent (pix intents are
		// cancellable; cards already charged can't be). The purchase stays
		// pending and the expiry worker frees the reservation in 10:01.
		o.cancelIntentBestEffort(ctx, intent.ID)
		return nil, err
	}

	return o.result(res), nil
}

// ── Result ─────────────────────────────────────────────────────────────────

func (o *Operations) result(res *reservation) *CheckoutResult {
	return &CheckoutResult{
		PurchaseID:       res.purchase.ID,
		EditionID:        res.purchase.EditionID,
		Status:           res.purchase.Status,
		StatusReason:     res.purchase.StatusReason,
		TotalCents:       res.purchase.TotalCents,
		Currency:         res.purchase.Currency,
		PaymentMethod:    res.purchase.PaymentMethod,
		PayssageIntentID: res.purchase.PayssageIntentID,
		QRCode:           res.purchase.QRCode,
		QRCodeBase64:     res.purchase.QRCodeBase64,
		ExpiresAt:        res.purchase.ExpiresAt,
		CreatedAt:        res.purchase.CreatedAt,
		Items:            res.items,
		WsToken:          res.wsToken,
		WsTokenExpiresAt: res.wsTokenExp,
	}
}

// ── Internals ──────────────────────────────────────────────────────────────

type reservation struct {
	purchase   *models.Purchase
	items      []models.PurchaseItem
	priced     []pricedLine
	wsToken    string
	wsTokenExp time.Time
	free       bool
}

type pricedLine struct {
	Line      CheckoutLine
	UnitPrice int64
	Base      *int   // the item's base quantity (nil = unlimited), from the locked row
	Name      string // display name (ticket type / product variant / program) — sent to the provider as the item title
}

// outOfStockError is the 409 listing the items that cannot be fully
// satisfied (availability semantics from split 3).
func outOfStockError(itemIDs []uuid.UUID) error {
	fields := make([]fielderr.Field, 0, len(itemIDs))
	for _, id := range itemIDs {
		fields = append(fields, fielderr.Field{Field: "items", Message: "not enough stock for item " + id.String()})
	}
	return fielderr.Wrap(errs.ErrItemsUnavailable, fields...)
}

func itemNotInEditionError(line CheckoutLine) error {
	return fielderr.On(errs.ErrItemNotInEdition, fmt.Sprintf("items[].%s", line.ItemType), line.ItemID.String())
}

// mapItemLoadError turns a not-found from the FOR UPDATE reads into the 400
// item-list contract (unknown item ids are validation, not 404 — the cart
// the buyer submitted is wrong, the endpoint exists).
func (o *Operations) mapItemLoadError(err error, line CheckoutLine) error {
	if errors.Is(err, database.ErrNotFound) {
		return fielderr.On(errs.ErrUnknownItem, fmt.Sprintf("items[].%s", line.ItemType), line.ItemID.String())
	}
	return err
}

// giftEmailDelay holds the gift email back long enough for the confirming
// tx to commit: it is enqueued inside the tx, before commit.
const giftEmailDelay = 30 * time.Second

// enqueueGiftEmail schedules the gifted-ticket email for an accountless
// free-order recipient inside the checkout tx (the registration is already
// confirmed at checkout). The delay keeps the job from running before the
// confirmation commits; on rollback it finds no registration and no-ops.
// Account holders are skipped: they see the ticket in my-ticket and get
// the badge email.
func (o *Operations) enqueueGiftEmail(ctx context.Context, reg *models.Registration) {
	if reg.AttendeeUserID != nil {
		return
	}
	err := o.Jobs.EnqueueAfter(ctx, jobs.SendGiftEmailArgs{RegistrationID: reg.ID}, giftEmailDelay)
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "checkout: failed to enqueue gift email",
			slog.String("registration_id", reg.ID.String()),
			telemetry.Err(err))
	}
}

// enqueueExpiry enqueues the purchases.expire job inside the checkout tx,
// before commit, delayed to the purchase's expires_at. The delay means it
// cannot run before the purchase is visible; an enqueue failure fails the
// checkout; a rolled-back checkout leaves a job that finds no purchase and
// no-ops.
func (o *Operations) enqueueExpiry(ctx context.Context, purchaseID uuid.UUID, expiresAt time.Time) error {
	err := o.Jobs.EnqueueAfter(ctx, jobs.ExpirePurchaseArgs{PurchaseID: purchaseID}, time.Until(expiresAt))
	if err != nil {
		return fmt.Errorf("checkout: enqueue expiry: %w", err)
	}
	return nil
}

// notifyStock publishes the reservation's stock deltas (D10): item_ids only
// — the SSE relay re-queries availability from the DB. Fire-and-forget: a
// missed notification is a stale snapshot, never data loss.
func (o *Operations) notifyStock(ctx context.Context, purchase *models.Purchase, items []models.PurchaseItem) {
	stock := notify.StockNotification{Kind: notify.KindStock, EditionID: purchase.EditionID}
	for _, item := range items {
		stock.ItemIDs = append(stock.ItemIDs, item.ItemID)
	}
	o.publish(ctx, stock)
}

// notifyPurchase publishes a purchase event (D9) — the free-order confirm
// and the revert cancel. The WS hub maps the status to the frame.
func (o *Operations) notifyPurchase(ctx context.Context, purchase *models.Purchase, status models.PurchaseStatus) {
	o.publish(ctx, notify.PurchaseNotification{
		Kind:       notify.KindPurchase,
		EditionID:  purchase.EditionID,
		PurchaseID: purchase.ID,
		Status:     string(status),
	})
}

func (o *Operations) publish(ctx context.Context, payload any) {
	raw, err := json.Marshal(payload)
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "checkout: marshal notifier payload",
			slog.String("channel", channelUniventsChanges),
			telemetry.Err(err))
		return
	}
	err = o.Notifier.Notify(ctx, channelUniventsChanges, string(raw))
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "checkout: publish to notifier",
			slog.String("channel", channelUniventsChanges),
			slog.String("payload", string(raw)),
			telemetry.Err(err))
	}
}

// revert cancels a purchase whose intent could not be created (rejected
// card, provider down, revoked seller…): guarded pending→cancelled,
// materialized rows flipped to cancelled (stock freed), NOTIFY. Best-effort:
// the original payment error is what the buyer sees; failures here are
// logged (the expiry worker is the backstop).
func (o *Operations) revert(ctx context.Context, purchase *models.Purchase, items []models.PurchaseItem) {
	err := o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.Purchases.UpdateStatusIf(ctx, purchase.ID,
			models.PurchaseStatusPending, models.PurchaseStatusCancelled, nil)
		if err != nil {
			return err
		}
		if updated == nil {
			return nil // already flipped (e.g. a webhook won the race)
		}
		return o.flipCancelled(ctx, items)
	})
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "checkout: revert failed after intent creation error",
			slog.String("purchase_id", purchase.ID.String()),
			telemetry.Err(err))
		return
	}
	o.notifyStock(ctx, purchase, items)
	o.notifyPurchase(ctx, purchase, models.PurchaseStatusCancelled)
}

// flipCancelled materializes the revert (D4), inside the revert tx:
// registrations/product_purchases/participations → cancelled.
func (o *Operations) flipCancelled(ctx context.Context, items []models.PurchaseItem) error {
	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				continue
			}
			_, err := o.Registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				continue
			}
			_, err := o.ProductPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				continue
			}
			_, err := o.Participations.UpdateParticipationStatus(ctx, *item.ParticipationID, models.ProgramParticipationStatusCancelled)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
