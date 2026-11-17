package service

import (
	"context"
	"time"
	"uuid"

	"univents/errs"

	"univents/models"
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
