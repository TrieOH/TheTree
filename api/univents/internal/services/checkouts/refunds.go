package checkouts

import (
	"context"

	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// PurchaseAttendee is one person a ticket line was assigned to — surfaced on
// the organizer orders read so the admin can see who (besides the payer) the
// purchase covers. Gifted tickets show the recipient; the payer_email on the
// purchase identifies who gets the refund.
type PurchaseAttendee struct {
	Name  string
	Email string
}

// EditionPurchase is the organizer orders read row: the shared purchase
// shape plus its items and the ticket attendees (from registrations).
type EditionPurchase struct {
	Purchase  models.Purchase
	Items     []models.PurchaseItem
	Attendees []PurchaseAttendee
}

// ListEditionPurchases is the organizer orders read (refund plan B3): every
// purchase of an edition, newest first, with items and attendee names.
// Owner/admin only — non-members 404, staff 403 (via CheckEvent).
func (o *Operations) ListEditionPurchases(ctx context.Context, actorID, editionID uuid.UUID) ([]EditionPurchase, error) {
	edition, err := o.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}
	err = o.authz.CheckEvent(ctx, actorID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	purchases, err := o.purchases.ListByEdition(ctx, editionID)
	if err != nil {
		return nil, err
	}
	out := make([]EditionPurchase, 0, len(purchases))
	for _, p := range purchases {
		ep := EditionPurchase{Purchase: p}
		ep.Items, err = o.purchases.ListItemsByPurchase(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		ep.Attendees, err = o.attendeesFor(ctx, ep.Items)
		if err != nil {
			return nil, err
		}
		out = append(out, ep)
	}
	return out, nil
}

// attendeesFor resolves the ticket items' registrations to attendee
// name/email (deduped — a purchase can have several ticket lines).
func (o *Operations) attendeesFor(ctx context.Context, items []models.PurchaseItem) ([]PurchaseAttendee, error) {
	var out []PurchaseAttendee
	seen := make(map[uuid.UUID]bool)
	for _, item := range items {
		if item.ItemType != models.PurchaseItemTypeTicket || item.RegistrationID == nil {
			continue
		}
		if seen[*item.RegistrationID] {
			continue
		}
		seen[*item.RegistrationID] = true
		reg, err := o.registrations.GetByID(ctx, *item.RegistrationID)
		if err != nil {
			return nil, err
		}
		out = append(out, PurchaseAttendee{Name: reg.AttendeeName, Email: reg.AttendeeEmail})
	}
	return out, nil
}

// RefundPurchase initiates a full refund of an approved purchase (refund
// plan B3): owner/admin only. It calls payssage RefundIntent and returns the
// purchase — which stays `approved` until the payment.refunded webhook flips
// it (webhook-confirmed, single writer, D-2). An already-refunded purchase
// is a 409; a pending/expired/cancelled one cannot be refunded.
func (o *Operations) RefundPurchase(ctx context.Context, actorID, purchaseID uuid.UUID) (*models.Purchase, error) {
	purchase, err := o.purchases.GetByID(ctx, purchaseID)
	if err != nil {
		return nil, err
	}
	edition, err := o.editions.GetByID(ctx, purchase.EditionID)
	if err != nil {
		return nil, err
	}
	err = o.authz.CheckEvent(ctx, actorID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	switch purchase.Status {
	case models.PurchaseStatusApproved:
		// ok
	case models.PurchaseStatusRefunded:
		return nil, fun.Err("purchase already refunded").Conflict()
	default:
		return nil, fun.Errf("purchase cannot be refunded from status %q", purchase.Status).Conflict()
	}
	if purchase.PayssageIntentID == nil {
		return nil, fun.Err("purchase has no payment intent (free order) and cannot be refunded").BadRequest()
	}

	_, err = o.payssage.RefundIntent(ctx, *purchase.PayssageIntentID)
	if err != nil {
		return nil, mapPayssageError(err)
	}
	// The webhook flips the purchase to refunded; return the current state.
	return purchase, nil
}
