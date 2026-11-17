package service

import (
	"context"
	"errors"
	"lib/database"
	"lib/fielderr"
	"sort"
	"time"
	"univents/errs"
	"univents/models"
	"uuid"
)

// reserve is the checkout transaction: locks the item rows in a
// deterministic order (no deadlocks), recomputes prices from the DB,
// enforces availability, and materializes the purchase — the purchase row,
// its purchase_items, the materialized registrations / product_purchases /
// program_participations, the delayed expiry job (paid orders), and the
// one-time ws token — all in one tx. Free orders branch inside: approved +
// rows confirmed (+ badge emit), no job.
func (o *Operations) reserve(ctx context.Context, editionID, purchaserID uuid.UUID, lines []CheckoutLine, in CheckoutInput) (*reservation, error) {
	expiresAt := time.Now().UTC().Add(reservationDuration)
	var res reservation

	err := o.Tx.WithinTx(ctx, func(ctx context.Context) error {
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

		created, err := o.Purchases.CreatePurchase(ctx, purchase)
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
			err = o.enqueueExpiry(ctx, created.ID, expiresAt)
			if err != nil {
				return err
			}
		}

		raw, tokenExp, err := o.Tokens.IssueToken(ctx, created.ID, purchaserID)
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
			tt, err := o.TicketTypes.GetByIDForUpdate(ctx, line.ItemID)
			if err != nil {
				return nil, o.mapItemLoadError(err, line)
			}
			if tt.EditionID != editionID {
				return nil, itemNotInEditionError(line)
			}
			out = append(out, pricedLine{Line: line, UnitPrice: tt.PriceCents, Base: tt.MaxQuantity, Name: tt.Name})

		case models.PurchaseItemTypeProduct:
			v, err := o.Products.GetVariantByIDForUpdate(ctx, line.ItemID)
			if err != nil {
				return nil, o.mapItemLoadError(err, line)
			}
			if v.EditionID != editionID {
				return nil, itemNotInEditionError(line)
			}
			out = append(out, pricedLine{Line: line, UnitPrice: v.Price, Base: v.Stock, Name: v.Name})

		case models.PurchaseItemTypeProgramOccurrence:
			occ, err := o.Occurrences.GetOccurrenceByIDForUpdate(ctx, line.ItemID)
			if err != nil {
				return nil, o.mapItemLoadError(err, line)
			}
			if occ.EditionID != editionID {
				return nil, itemNotInEditionError(line)
			}
			prog, err := o.Programs.GetByIDForUpdate(ctx, occ.ProgramID)
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
			return nil, fielderr.On(errs.ErrInvalidItemType, "items[].item_type", string(line.ItemType))
		}
	}
	return out, nil
}

// checkAvailability reads the edition's stock position inside the tx (the
// availability ledger sums purchase_items of pending/approved purchases;
// the current purchase is not inserted yet, so it never counts itself) and
// reports the items that cannot be fully satisfied.
func (o *Operations) checkAvailability(ctx context.Context, editionID uuid.UUID, priced []pricedLine) ([]uuid.UUID, error) {
	avail, err := o.Purchases.Availability(ctx, editionID)
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
		reg, err := o.Registrations.Create(ctx, &models.Registration{
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
			_, err := o.Badges.EmitForConfirmedRegistration(ctx, reg.ID)
			if err != nil {
				return nil, err
			}
			// Email-only free gift: schedule the gifted-ticket email in the
			// same tx (accountless recipients get the claim instructions;
			// the badge itself is deferred until they claim).
			o.enqueueGiftEmail(ctx, reg)
		}
		item, err := o.Purchases.CreatePurchaseItem(ctx, &models.PurchaseItem{
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
			pp, err := o.ProductPurchases.CreateProductPurchase(ctx, &models.ProductPurchase{
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
			item, err := o.Purchases.CreatePurchaseItem(ctx, &models.PurchaseItem{
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
				return nil, errs.ErrProgramNeedsTicket
			}
			part, err := o.Participations.CreateParticipation(ctx, &models.ProgramParticipation{
				EditionID:      purchase.EditionID,
				OccurrenceID:   pl.Line.ItemID,
				RegistrationID: *firstTicketRegID,
				Status:         models.ProgramParticipationStatusRegistered,
			})
			if err != nil {
				return nil, err
			}
			item, err := o.Purchases.CreatePurchaseItem(ctx, &models.PurchaseItem{
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
// edition, or nil when they hold none. The repo returns an absent row as
// database.ErrNotFound; here, "no ticket" is a normal state — the one-ticket check
// and the my-ticket read both fall back to it.
func (o *Operations) activeRegistrationFor(ctx context.Context, editionID, attendeeID uuid.UUID) (*models.Registration, error) {
	reg, err := o.Registrations.GetActiveByEditionAndAttendee(ctx, editionID, attendeeID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
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
		existing, err := o.Registrations.GetActiveByEditionAndAttendeeEmail(ctx, editionID, pl.Line.Attendee.Email)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
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
	return fielderr.On(errs.ErrPersonHasTicket, field, message)
}
