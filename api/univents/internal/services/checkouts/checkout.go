package checkouts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"lib/database"
	"lib/telemetry"

	"univents/internal/services/checkouts/jobs"
	"univents/internal/services/notify"
	"univents/models"
	"univents/ports"

	payssage "sdk/payssage"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"go.uber.org/zap"
)

// reservationDuration is how long an unpaid reservation is held: 10:01 (the
// master plan's expiry budget). The purchase's expires_at and the expiry
// job's ScheduledAt share this one value.
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
	ctx, span := telemetry.StartSpan(ctx, "CheckoutsService.Checkout")
	defer span.End()

	// 1. Edition + event + payment config (plain reads — the tx only locks
	//    the item rows).
	edition, err := o.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}
	if edition.IsDraft {
		return nil, fun.ErrBadRequest("edition is not published")
	}
	event, err := o.events.GetByID(ctx, edition.EventID)
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
	//    expiry job (river.InsertTx), issue the ws token — commit.
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
		return nil, fun.ErrBadRequest("event has no payment config")
	}
	err = o.validatePayment(in)
	if err != nil {
		o.revert(ctx, res.purchase, res.items)
		return nil, err
	}

	// 4. Create the Payssage intent (post-commit — no HTTP inside the DB
	//    tx; if this fails the purchase is reverted so the buyer can retry
	//    immediately).
	intent, err := o.payssage.Checkout(ctx, o.walletID, o.buildIntentRequest(res.purchase, in, *event.PayssageSellerID, res.priced))
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

// reserve is the checkout transaction: locks the item rows in a
// deterministic order (no deadlocks), recomputes prices from the DB,
// enforces availability, and materializes the purchase — the purchase row,
// its purchase_items, the materialized registrations / product_purchases /
// program_participations, the expiry river job (paid orders), and the
// one-time ws token — all in one tx. Free orders branch inside: approved +
// rows confirmed (+ badge emit), no job.
func (o *Operations) reserve(ctx context.Context, editionID, purchaserID uuid.UUID, lines []CheckoutLine, in CheckoutInput) (*reservation, error) {
	expiresAt := time.Now().UTC().Add(reservationDuration)
	var res reservation

	err := o.tx.WithinTx(ctx, func(ctx context.Context) error {
		priced, err := o.lockAndPrice(ctx, editionID, lines)
		if err != nil {
			return err
		}

		unavailable, err := o.checkAvailability(ctx, editionID, priced)
		if err != nil {
			return err
		}
		if len(unavailable) > 0 {
			return outOfStockError(unavailable)
		}

		// One active ticket per person per edition: reject the cart when any
		// ticket line's attendee already holds a pending/confirmed
		// registration here. The partial unique index
		// (uniq_registrations_active_per_edition_attendee) is the backstop
		// for concurrent checkouts on the same attendee.
		err = o.checkOneTicketPerPerson(ctx, editionID, priced)
		if err != nil {
			return err
		}

		total := int64(0)
		for _, pl := range priced {
			total += int64(pl.Line.Quantity) * pl.UnitPrice
		}
		free := total == 0

		purchase := &models.Purchase{
			EditionID:     editionID,
			PurchaserID:   purchaserID,
			Status:        models.PurchaseStatusPending,
			TotalCents:    total,
			Currency:      currencyBRL,
			PaymentMethod: &in.PaymentMethod,
			ExpiresAt:     expiresAt,
		}
		// Payer identity for refunds (refund plan B3): who the provider will
		// refund. The front sends payment?.payer_email ?? profile.email.
		if in.Payer.Email != "" {
			purchase.PayerEmail = new(in.Payer.Email)
		}
		if free {
			purchase.Status = models.PurchaseStatusApproved
			purchase.PaymentMethod = nil // no payment happened
		}

		created, err := o.purchases.CreatePurchase(ctx, purchase)
		if err != nil {
			return err // unique pending-per-purchaser-edition → 409 via the registry message
		}
		res.purchase = created

		items, err := o.materialize(ctx, created, priced, free)
		if err != nil {
			return err
		}
		res.items = items

		if !free {
			job, err := o.insertExpiryJob(ctx, created.ID, expiresAt)
			if err != nil {
				return err
			}
			updated, err := o.purchases.UpdateRiverJobID(ctx, created.ID, job.Job.ID)
			if err != nil {
				return err
			}
			res.purchase = updated
		}

		raw, tokenExp, err := o.tokens.IssueToken(ctx, created.ID, purchaserID)
		if err != nil {
			return err
		}
		res.wsToken = raw
		res.wsTokenExp = tokenExp
		res.free = free
		res.priced = priced
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// lockAndPrice loads every requested item row FOR UPDATE in a deterministic
// order (item type, then id — two checkouts locking in different orders
// could deadlock), verifies the item belongs to the edition, and computes
// the server-side unit price. The locks serialize concurrent checkouts on
// the same items before availability is checked — no oversell on the last
// unit.
func (o *Operations) lockAndPrice(ctx context.Context, editionID uuid.UUID, lines []CheckoutLine) ([]pricedLine, error) {
	sorted := append([]CheckoutLine(nil), lines...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].ItemType != sorted[j].ItemType {
			return sorted[i].ItemType < sorted[j].ItemType
		}
		return sorted[i].ItemID.String() < sorted[j].ItemID.String()
	})

	out := make([]pricedLine, 0, len(sorted))
	for _, line := range sorted {
		switch line.ItemType {
		case models.PurchaseItemTypeTicket:
			tt, err := o.ticketTypes.GetByIDForUpdate(ctx, line.ItemID)
			if err != nil {
				return nil, o.mapItemLoadError(err, line)
			}
			if tt.EditionID != editionID {
				return nil, itemNotInEditionError(line)
			}
			out = append(out, pricedLine{Line: line, UnitPrice: tt.PriceCents, Base: tt.MaxQuantity, Name: tt.Name})

		case models.PurchaseItemTypeProduct:
			v, err := o.products.GetVariantByIDForUpdate(ctx, line.ItemID)
			if err != nil {
				return nil, o.mapItemLoadError(err, line)
			}
			if v.EditionID != editionID {
				return nil, itemNotInEditionError(line)
			}
			out = append(out, pricedLine{Line: line, UnitPrice: v.Price, Base: v.Stock, Name: v.Name})

		case models.PurchaseItemTypeProgramOccurrence:
			occ, err := o.occurrences.GetOccurrenceByIDForUpdate(ctx, line.ItemID)
			if err != nil {
				return nil, o.mapItemLoadError(err, line)
			}
			if occ.EditionID != editionID {
				return nil, itemNotInEditionError(line)
			}
			prog, err := o.programs.GetByIDForUpdate(ctx, occ.ProgramID)
			if err != nil {
				return nil, o.mapItemLoadError(err, line)
			}
			if prog.EditionID != editionID {
				return nil, itemNotInEditionError(line)
			}
			price := int64(0) // programs.price is nullable — nil = free
			if prog.Price != nil {
				price = *prog.Price
			}
			out = append(out, pricedLine{Line: line, UnitPrice: price, Base: occ.MaxCapacity, Name: prog.Name})

		default:
			return nil, fun.Err("invalid item_type").WithFields(&fun.FieldError{
				Field: "items[].item_type", Message: string(line.ItemType),
			}).BadRequest()
		}
	}
	return out, nil
}

// checkAvailability reads the edition's stock position inside the tx (the
// availability ledger sums purchase_items of pending/approved purchases;
// the current purchase is not inserted yet, so it never counts itself) and
// reports the items that cannot be fully satisfied.
func (o *Operations) checkAvailability(ctx context.Context, editionID uuid.UUID, priced []pricedLine) ([]uuid.UUID, error) {
	avail, err := o.purchases.Availability(ctx, editionID)
	if err != nil {
		return nil, err
	}
	stock := make(map[uuid.UUID]models.ItemAvailability, len(avail))
	for _, a := range avail {
		stock[a.ItemID] = a
	}

	var unavailable []uuid.UUID
	for _, pl := range priced {
		a, ok := stock[pl.Line.ItemID]
		if !ok {
			unavailable = append(unavailable, pl.Line.ItemID)
			continue
		}
		if a.BaseQuantity == nil {
			continue // unlimited — never sold out
		}
		if int64(*a.BaseQuantity)-a.ReservedQuantity < int64(pl.Line.Quantity) {
			unavailable = append(unavailable, pl.Line.ItemID)
		}
	}
	return unavailable, nil
}

// materialize creates the purchase_items rows and the materialized rows
// (D4), inside the checkout tx: ticket units → one registration each;
// product lines → one product_purchase (attached to the first ticket's
// registration when a ticket exists); program lines → one participation on
// the first ticket's registration. Tickets are created FIRST regardless of
// the lock order (lockAndPrice sorts for locking; materialization needs the
// ticket registration before products/programs attach to it). Free orders
// create rows directly in their confirmed state (+ badge emit); paid orders
// create them pending — the webhook receiver flips them on approval.
func (o *Operations) materialize(ctx context.Context, purchase *models.Purchase, priced []pricedLine, free bool) ([]models.PurchaseItem, error) {
	items := make([]models.PurchaseItem, 0, len(priced))
	var firstTicketRegID *uuid.UUID

	// Pass 1: tickets (also seeds firstTicketRegID for the passes below).
	for _, pl := range priced {
		if pl.Line.ItemType != models.PurchaseItemTypeTicket {
			continue
		}
		status := models.RegistrationStatusPending
		if free {
			status = models.RegistrationStatusConfirmed
		}
		reg, err := o.registrations.Create(ctx, &models.Registration{
			EditionID:      purchase.EditionID,
			TicketTypeID:   pl.Line.ItemID,
			PurchaserID:    purchase.PurchaserID,
			AttendeeUserID: pl.Line.Attendee.UserID,
			AttendeeEmail:  pl.Line.Attendee.Email,
			AttendeeName:   pl.Line.Attendee.Name,
			Status:         status,
		})
		if err != nil {
			return nil, err
		}
		if firstTicketRegID == nil {
			firstTicketRegID = &reg.ID
		}
		if free {
			// A free ticket is confirmed at checkout — emit its badge.
			_, err := o.badges.EmitForConfirmedRegistration(ctx, reg.ID)
			if err != nil {
				return nil, err
			}
			// Email-only free gift: schedule the gifted-ticket email in the
			// same tx (accountless recipients get the claim instructions;
			// the badge itself is deferred until they claim).
			o.enqueueGiftEmail(ctx, reg)
		}
		item, err := o.purchases.CreatePurchaseItem(ctx, &models.PurchaseItem{
			PurchaseID:     purchase.ID,
			ItemType:       models.PurchaseItemTypeTicket,
			ItemID:         pl.Line.ItemID,
			Quantity:       1,
			UnitPriceCents: pl.UnitPrice,
			RegistrationID: &reg.ID,
		})
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	// Pass 2: products + programs.
	for _, pl := range priced {
		switch pl.Line.ItemType {
		case models.PurchaseItemTypeProduct:
			status := models.ProductPurchaseStatusPending
			if free {
				status = models.ProductPurchaseStatusConfirmed
			}
			pp, err := o.productPurchases.CreateProductPurchase(ctx, &models.ProductPurchase{
				EditionID:      purchase.EditionID,
				VariantID:      pl.Line.ItemID,
				PurchaserID:    purchase.PurchaserID,
				RegistrationID: firstTicketRegID, // attach to the ticket's registration when present
				Quantity:       pl.Line.Quantity,
				Status:         status,
			})
			if err != nil {
				return nil, err
			}
			item, err := o.purchases.CreatePurchaseItem(ctx, &models.PurchaseItem{
				PurchaseID:        purchase.ID,
				ItemType:          models.PurchaseItemTypeProduct,
				ItemID:            pl.Line.ItemID,
				Quantity:          pl.Line.Quantity,
				UnitPriceCents:    pl.UnitPrice,
				ProductPurchaseID: &pp.ID,
			})
			if err != nil {
				return nil, err
			}
			items = append(items, *item)

		case models.PurchaseItemTypeProgramOccurrence:
			if firstTicketRegID == nil {
				// Unreachable — validateLines enforces ≥1 ticket whenever a
				// program line is present (D4b).
				return nil, fun.ErrBadRequest("program items require a ticket item")
			}
			part, err := o.participations.CreateParticipation(ctx, &models.ProgramParticipation{
				EditionID:      purchase.EditionID,
				OccurrenceID:   pl.Line.ItemID,
				RegistrationID: *firstTicketRegID,
				Status:         models.ProgramParticipationStatusRegistered,
			})
			if err != nil {
				return nil, err
			}
			item, err := o.purchases.CreatePurchaseItem(ctx, &models.PurchaseItem{
				PurchaseID:      purchase.ID,
				ItemType:        models.PurchaseItemTypeProgramOccurrence,
				ItemID:          pl.Line.ItemID,
				Quantity:        1,
				UnitPriceCents:  pl.UnitPrice,
				ParticipationID: &part.ID,
			})
			if err != nil {
				return nil, err
			}
			items = append(items, *item)
		}
	}
	return items, nil
}

// activeRegistrationFor returns the attendee's active registration in the
// edition, or nil when they hold none. The repo surfaces an absent row as
// NOT_FOUND; here, "no ticket" is a normal state — the one-ticket check
// and the my-ticket read both fall back to it.
func (o *Operations) activeRegistrationFor(ctx context.Context, editionID, attendeeID uuid.UUID) (*models.Registration, error) {
	reg, err := o.registrations.GetActiveByEditionAndAttendee(ctx, editionID, attendeeID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			//nolint:nilnil // holds no ticket is a normal state, not an error
			return nil, nil
		}
		return nil, err
	}
	return reg, nil
}

// checkOneTicketPerPerson enforces the one-ticket-per-person rule inside
// the checkout tx: every ticket line's attendee must not already hold an
// active (pending or confirmed) registration in the edition — for
// themselves or as a gifted recipient. Attendees with an account are
// checked by user id; accountless (email-only) recipients by email. The
// clean 409 comes from this pre-check; the partial unique indexes
// (uniq_registrations_active_per_edition_attendee / _email_per_edition)
// are the backstop that catch two concurrent checkouts racing on the same
// attendee atomically.
func (o *Operations) checkOneTicketPerPerson(ctx context.Context, editionID uuid.UUID, priced []pricedLine) error {
	for _, pl := range priced {
		if pl.Line.ItemType != models.PurchaseItemTypeTicket {
			continue
		}
		if pl.Line.Attendee.UserID != nil {
			existing, err := o.activeRegistrationFor(ctx, editionID, *pl.Line.Attendee.UserID)
			if err != nil {
				return err
			}
			if existing != nil {
				return oneTicketPerPersonError("items[].ticket.attendee.user_id", pl.Line.Attendee.UserID.String())
			}
			continue
		}
		// Email-only recipient (no IdentityX account yet): the slot is held
		// by email — another email-only gift to the same address is blocked
		// (the partial unique index is the concurrency backstop). Absent
		// row = holds no ticket, a normal state.
		existing, err := o.registrations.GetActiveByEditionAndAttendeeEmail(ctx, editionID, pl.Line.Attendee.Email)
		if err != nil {
			if fun.Is(err, fun.CodeNotFound) {
				continue
			}
			return err
		}
		if existing != nil {
			return oneTicketPerPersonError("items[].ticket.attendee.email", pl.Line.Attendee.Email)
		}
	}
	return nil
}

// oneTicketPerPersonError is the clean 409 for a cart/DB violation of the
// one-ticket-per-person rule.
func oneTicketPerPersonError(field, message string) error {
	return fun.Err("this person already has a ticket for this edition").WithFields(&fun.FieldError{
		Field: field, Message: message,
	}).Conflict()
}

// enqueueGiftEmail schedules the gifted-ticket email for an accountless
// free-order recipient inside the checkout tx (the registration is already
// confirmed at checkout — the gift email and the confirmation commit
// atomically). Account holders are skipped: they see the ticket in
// my-ticket and get the badge email.
func (o *Operations) enqueueGiftEmail(ctx context.Context, reg *models.Registration) {
	if reg.AttendeeUserID != nil {
		return
	}
	tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx)
	if !ok {
		telemetry.Log().Error("checkout: no transaction in context for gift email",
			zap.String("registration_id", reg.ID.String()))
		return
	}
	_, err := o.river.InsertTx(ctx, tx, jobs.SendGiftEmailArgs{RegistrationID: reg.ID}, nil)
	if err != nil {
		telemetry.Log().Error("checkout: failed to enqueue gift email",
			zap.String("registration_id", reg.ID.String()),
			zap.Error(err))
	}
}

// insertExpiryJob enqueues the purchases.expire job inside the checkout tx
// (river.InsertTx against the tx from context — the payssage
// dispatchDeliveries pattern), so the job + purchase commit atomically.
func (o *Operations) insertExpiryJob(ctx context.Context, purchaseID uuid.UUID, scheduledAt time.Time) (*rivertype.JobInsertResult, error) {
	tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx)
	if !ok {
		return nil, fun.Err("checkout: no transaction in context").Internal()
	}
	return o.river.InsertTx(ctx, tx, jobs.ExpirePurchaseArgs{PurchaseID: purchaseID}, &river.InsertOpts{ScheduledAt: scheduledAt})
}

// attachIntent stores the intent on the purchase in a second, short
// transaction (post-commit): seller, intent id (the D2 correlation key),
// and the pix QR from the intent's provider_data.
func (o *Operations) attachIntent(ctx context.Context, res *reservation, intent *payssage.Intent) error {
	qr, qrB64 := pixQR(intent.ProviderData)
	return o.tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.purchases.AttachIntent(ctx, res.purchase.ID, intent.SellerID, intent.ID, qr, qrB64)
		if err != nil {
			return err
		}
		res.purchase = updated
		return nil
	})
}

// revert cancels a purchase whose intent could not be created (rejected
// card, provider down, revoked seller…): guarded pending→cancelled,
// materialized rows flipped to cancelled (stock freed), NOTIFY. Best-effort:
// the original payment error is what the buyer sees; failures here are
// logged (the expiry worker is the backstop).
func (o *Operations) revert(ctx context.Context, purchase *models.Purchase, items []models.PurchaseItem) {
	ctx, span := telemetry.StartSpan(ctx, "CheckoutsService.Revert")
	defer span.End()
	err := o.tx.WithinTx(ctx, func(ctx context.Context) error {
		updated, err := o.purchases.UpdateStatusIf(ctx, purchase.ID,
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
		telemetry.Log().Error("checkout: revert failed after intent creation error",
			zap.String("purchase_id", purchase.ID.String()),
			zap.Error(err))
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
			_, err := o.registrations.UpdateStatus(ctx, *item.RegistrationID, models.RegistrationStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				continue
			}
			_, err := o.productPurchases.UpdateProductPurchaseStatus(ctx, *item.ProductPurchaseID, models.ProductPurchaseStatusCancelled, nil)
			if err != nil {
				return err
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				continue
			}
			_, err := o.participations.UpdateParticipationStatus(ctx, *item.ParticipationID, models.ProgramParticipationStatusCancelled)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// cancelIntentBestEffort cancels a stray intent (Tx 2 failure): pix intents
// are cancellable; a card already charged comes back non-cancellable and is
// left for the refund follow-up (D1). Log-only.
func (o *Operations) cancelIntentBestEffort(ctx context.Context, intentID uuid.UUID) {
	_, err := o.payssage.CancelIntent(ctx, intentID)
	if err != nil {
		telemetry.Log().Warn("checkout: best-effort intent cancel failed",
			zap.String("intent_id", intentID.String()),
			zap.Error(err))
	}
}

// ── Validation helpers ─────────────────────────────────────────────────────

func (o *Operations) validatePayment(in CheckoutInput) error {
	switch in.PaymentMethod {
	case paymentPix:
		return nil
	case paymentCreditCard, paymentDebitCard:
		if in.CardToken == nil || *in.CardToken == "" {
			return fun.Err("card_token is required for card payments").WithFields(&fun.FieldError{
				Field: "card_token", Message: "required",
			}).BadRequest()
		}
		if in.PaymentMethodID == nil || *in.PaymentMethodID == "" {
			return fun.Err("payment_method_id is required for card payments").WithFields(&fun.FieldError{
				Field: "payment_method_id", Message: "required",
			}).BadRequest()
		}
		if in.Installments != nil && *in.Installments < 1 {
			return fun.Err("installments must be at least 1").WithFields(&fun.FieldError{
				Field: "installments", Message: "min=1",
			}).BadRequest()
		}
		return nil
	default:
		return fun.Err("invalid payment_method").WithFields(&fun.FieldError{
			Field: "payment_method", Message: "one of pix, credit_card, debit_card",
		}).BadRequest()
	}
}

// validateLines enforces the cart shape: quantity > 0; ticket lines are
// one-per-person (quantity 1 + an attendee — gifting sends one line per
// person, so the same ticket type may repeat); product/program lines are
// unique per item; program lines are quantity 1 and require ≥1 ticket line
// in the same cart (D4b — a participation attaches to the ticket's
// registration).
func (o *Operations) validateLines(items []CheckoutLine) ([]CheckoutLine, error) {
	if len(items) == 0 {
		return nil, fun.Err("checkout requires at least one item").WithFields(&fun.FieldError{
			Field: "items", Message: "min=1",
		}).BadRequest()
	}

	seenProduct := make(map[uuid.UUID]bool, len(items))
	seenProgram := make(map[uuid.UUID]bool, len(items))
	seenAttendee := make(map[string]bool, len(items))
	hasTicket := false

	for i, line := range items {
		field := fmt.Sprintf("items[%d]", i)
		if line.Quantity < 1 {
			return nil, fun.Err("item quantity must be at least 1").WithFields(&fun.FieldError{
				Field: field + ".quantity", Message: "min=1",
			}).BadRequest()
		}
		switch line.ItemType {
		case models.PurchaseItemTypeTicket:
			if line.Quantity != 1 {
				return nil, fun.Err("ticket lines are one per person (quantity must be 1)").WithFields(&fun.FieldError{
					Field: field + ".quantity", Message: "tickets are one per person — send one line per attendee",
				}).BadRequest()
			}
			if line.Attendee == nil {
				return nil, fun.Err("ticket lines require an attendee").WithFields(&fun.FieldError{
					Field: field + ".attendee", Message: "required",
				}).BadRequest()
			}
			// One ticket per person per edition — two lines for the same
			// attendee (e.g. two tickets for yourself, or two gifts to the
			// same email) are rejected here; a ticket the attendee already
			// holds is caught in the tx (checkOneTicketPerPerson). The
			// email key is the pre-resolution pass — resolveAttendees
			// re-dedupes by the resolved id right after.
			attendeeKey := "e:" + strings.ToLower(strings.TrimSpace(line.Attendee.Email))
			if line.Attendee.UserID != nil {
				attendeeKey = "u:" + line.Attendee.UserID.String()
			}
			if seenAttendee[attendeeKey] {
				return nil, fun.Err("only one ticket per person").WithFields(&fun.FieldError{
					Field: field + ".attendee", Message: "this person already has a ticket in the cart",
				}).BadRequest()
			}
			seenAttendee[attendeeKey] = true
			hasTicket = true
		case models.PurchaseItemTypeProduct:
			if seenProduct[line.ItemID] {
				return nil, fun.Err("duplicate product line").WithFields(&fun.FieldError{
					Field: field, Message: "this product is already in the cart",
				}).BadRequest()
			}
			seenProduct[line.ItemID] = true
		case models.PurchaseItemTypeProgramOccurrence:
			if line.Quantity != 1 {
				return nil, fun.Err("program lines are one session per purchase (quantity must be 1)").WithFields(&fun.FieldError{
					Field: field + ".quantity", Message: "program lines are quantity 1",
				}).BadRequest()
			}
			if seenProgram[line.ItemID] {
				return nil, fun.Err("duplicate program line").WithFields(&fun.FieldError{
					Field: field, Message: "this program occurrence is already in the cart",
				}).BadRequest()
			}
			seenProgram[line.ItemID] = true
		default:
			return nil, fun.Err("invalid item_type").WithFields(&fun.FieldError{
				Field: field + ".item_type", Message: "one of ticket, product, program_occurrence",
			}).BadRequest()
		}
	}

	if len(seenProgram) > 0 && !hasTicket {
		return nil, fun.Err("program items require a ticket item in the same cart").WithFields(&fun.FieldError{
			Field: "items", Message: "a participation attaches to a ticket's registration",
		}).BadRequest()
	}
	return items, nil
}

// resolveAttendees resolves every ticket attendee's identity against
// IdentityX (HTTP — outside the checkout tx) and dedupes the cart:
//
//   - attendee with only an email → tie the account's actor id when the
//     email has one (existing account — the gift is already bound to the
//     person); otherwise the gift stays email-only (nil user id, the
//     recipient claims it after creating an account).
//   - attendee with both id and email → the email must resolve to that
//     same actor (GetByEmail + id compare): a different actor, or an email
//     with no account at all, is a 400 — the pair is inconsistent.
//   - two lines resolving to the same person (same id, or same email when
//     accountless) → 400, the one-ticket-per-person cart rule.
//
// Emails are normalized (trim + lowercase) to match IdentityX, which stores
// them lowercased. A failure to reach IdentityX fails the checkout closed:
// an attendee whose account cannot be verified must not be silently gifted
// as email-only (they may hold a ticket already).
func (o *Operations) resolveAttendees(ctx context.Context, lines []CheckoutLine) ([]CheckoutLine, error) {
	out := append([]CheckoutLine(nil), lines...)
	seen := make(map[string]bool, len(out))
	for i := range out {
		line := &out[i]
		if line.ItemType != models.PurchaseItemTypeTicket || line.Attendee == nil {
			continue
		}

		email := strings.TrimSpace(strings.ToLower(line.Attendee.Email))
		if email == "" {
			return nil, fun.Err("ticket lines require an attendee email").WithFields(&fun.FieldError{
				Field: fmt.Sprintf("items[%d].attendee.email", i), Message: "required",
			}).BadRequest()
		}
		line.Attendee.Email = email

		// One lookup per attendee: ErrActorNotFound means the email has no
		// account in the univents identityx project.
		actor, err := o.actors.GetByEmail(ctx, email)
		switch {
		case err != nil && !errors.Is(err, ports.ErrActorNotFound):
			// IdentityX unreachable — fail the checkout closed: an attendee
			// whose account cannot be verified must not be silently gifted
			// as email-only (they may already hold a ticket).
			return nil, err
		case errors.Is(err, ports.ErrActorNotFound) && line.Attendee.UserID != nil:
			// A user_id was claimed for an email with no account — the pair
			// cannot both be true.
			return nil, fun.Err("attendee email does not match an account").WithFields(&fun.FieldError{
				Field: fmt.Sprintf("items[%d].attendee.email", i), Message: email,
			}).BadRequest()
		case errors.Is(err, ports.ErrActorNotFound):
			// No account yet — email-only gift; keep the nil user id.
		case line.Attendee.UserID == nil:
			line.Attendee.UserID = &actor.ID
		case actor.ID != *line.Attendee.UserID:
			return nil, fun.Err("attendee user_id does not match the email").WithFields(&fun.FieldError{
				Field: fmt.Sprintf("items[%d].attendee.user_id", i), Message: line.Attendee.UserID.String(),
			}).BadRequest()
		}

		// Post-resolution dedup: two lines that resolve to the same person
		// (explicit id, or a resolved account) collide here; email-only
		// lines collide on the normalized email.
		key := "e:" + line.Attendee.Email
		if line.Attendee.UserID != nil {
			key = "u:" + line.Attendee.UserID.String()
		}
		if seen[key] {
			return nil, fun.Err("only one ticket per person").WithFields(&fun.FieldError{
				Field: fmt.Sprintf("items[%d].attendee", i), Message: "this person already has a ticket in the cart",
			}).BadRequest()
		}
		seen[key] = true
	}
	return out, nil
}

// ── Intent request ─────────────────────────────────────────────────────────

// buildIntentRequest maps the purchase + checkout input to the Payssage
// create-intent payload. For pix the server derives `payment_method_id` =
// "pix" (no tokenization exists); for cards the front's tokenization
// supplies it. `additional_info.items` carries the real product lines (name,
// quantity, unit price) so the provider's payment screen and risk engine
// see what the buyer is paying for — an intent without item data surfaces
// as a nameless product and can be flagged high-risk. Metadata carries
// purchase_id + edition_id for correlation.
func (o *Operations) buildIntentRequest(purchase *models.Purchase, in CheckoutInput, sellerID uuid.UUID, priced []pricedLine) payssage.CreateIntentRequest {
	providerData := map[string]any{
		"payer": map[string]any{
			"email":                 in.Payer.Email,
			"identification_type":   in.Payer.IdentificationType,
			"identification_number": in.Payer.IdentificationNumber,
		},
	}
	switch in.PaymentMethod {
	case paymentPix:
		providerData["payment_method_id"] = paymentPix
	case paymentCreditCard, paymentDebitCard:
		installments := 1
		if in.Installments != nil {
			installments = *in.Installments
		}
		providerData["payment_method_id"] = *in.PaymentMethodID
		providerData["token"] = *in.CardToken
		providerData["installments"] = installments
		if in.IssuerID != nil && *in.IssuerID != "" {
			providerData["issuer_id"] = *in.IssuerID
		}
	}

	if len(priced) > 0 {
		items := make([]map[string]any, 0, len(priced))
		for _, pl := range priced {
			items = append(items, map[string]any{
				"title":            pl.Name,
				"quantity":         pl.Line.Quantity,
				"unit_price_cents": pl.UnitPrice,
			})
		}
		providerData["additional_info"] = map[string]any{"items": items}
	}

	return payssage.CreateIntentRequest{
		SellerID:             sellerID,
		Currency:             currencyBRL,
		AmountCents:          purchase.TotalCents,
		CheckoutProviderData: providerData,
		Metadata: map[string]any{
			"purchase_id": purchase.ID.String(),
			"edition_id":  purchase.EditionID.String(),
		},
		// First-class correlation (refund plan A7): the caller ids payssage
		// stores as external_id (purchase) + external_group (edition). The
		// metadata copy is kept for backward compat.
		ExternalID:    new(purchase.ID.String()),
		ExternalGroup: new(purchase.EditionID.String()),
	}
}

// pixQR extracts the pix QR from the intent's provider_data
// (`pix_qr_code` / `pix_qr_code_base64` — payssage's MercadoPagoIntentData).
// Cards carry none — both stay nil.
func pixQR(providerData json.RawMessage) (qr, qrBase64 *string) {
	var d struct {
		PixQRCode       string `json:"pix_qr_code"`
		PixQRCodeBase64 string `json:"pix_qr_code_base64"`
	}
	if len(providerData) == 0 || json.Unmarshal(providerData, &d) != nil {
		return nil, nil
	}
	if d.PixQRCode != "" {
		qr = &d.PixQRCode
	}
	if d.PixQRCodeBase64 != "" {
		qrBase64 = &d.PixQRCodeBase64
	}
	return qr, qrBase64
}

// mapPayssageError converts a provider/API failure into a fun AppError the
// buyer can act on: 4xx (rejected card, invalid token...) → BAD_REQUEST with
// the provider's message; anything else stays an internal error. Without
// this the strict server would turn every payssage failure into a 500.
func mapPayssageError(err error) error {
	var apiErr *payssage.APIError
	if !errors.As(err, &apiErr) {
		return err
	}
	if apiErr.StatusCode >= 400 && apiErr.StatusCode < 500 {
		return fun.ErrBadRequest("payment failed: " + apiErr.Message)
	}
	return fun.ErrInternal("payment provider error: " + apiErr.Message)
}

// ── Errors ─────────────────────────────────────────────────────────────────

// outOfStockError is the 409 listing the items that cannot be fully
// satisfied (availability semantics from split 3).
func outOfStockError(itemIDs []uuid.UUID) error {
	fields := make([]any, 0, len(itemIDs))
	for _, id := range itemIDs {
		fields = append(fields, fun.FieldError{Field: "items", Message: "not enough stock for item " + id.String()})
	}
	return fun.Err("some items are no longer available").WithFields(fields...).Conflict()
}

func itemNotInEditionError(line CheckoutLine) error {
	return fun.Err("item does not belong to this edition").WithFields(&fun.FieldError{
		Field: fmt.Sprintf("items[].%s", line.ItemType), Message: line.ItemID.String(),
	}).BadRequest()
}

// mapItemLoadError turns a not-found from the FOR UPDATE reads into the 400
// item-list contract (unknown item ids are validation, not 404 — the cart
// the buyer submitted is wrong, the endpoint exists).
func (o *Operations) mapItemLoadError(err error, line CheckoutLine) error {
	if fun.Is(err, fun.CodeNotFound) {
		return fun.Err("unknown item").WithFields(&fun.FieldError{
			Field: fmt.Sprintf("items[].%s", line.ItemType), Message: line.ItemID.String(),
		}).BadRequest()
	}
	return err
}

// ── Notify ────────────────────────────────────────────────────────────────

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
		telemetry.Log().Error("checkout: marshal notifier payload",
			zap.String("channel", channelUniventsChanges),
			zap.Error(err))
		return
	}
	err = o.notifier.Notify(ctx, channelUniventsChanges, string(raw))
	if err != nil {
		telemetry.Log().Error("checkout: publish to notifier",
			zap.String("channel", channelUniventsChanges),
			zap.String("payload", string(raw)),
			zap.Error(err))
	}
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
