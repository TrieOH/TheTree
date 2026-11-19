package checkouts

import (
	"context"
	"errors"
	"fmt"
	"lib/fielderr"
	"strings"
	"univents/errs"
	"univents/models"
	"univents/ports"
	"uuid"
)

func (o *Operations) validatePayment(in CheckoutInput) error {
	switch in.PaymentMethod {
	case paymentPix:
		return nil
	case paymentCreditCard, paymentDebitCard:
		if in.CardToken == nil || *in.CardToken == "" {
			return fielderr.On(errs.ErrCardTokenRequired, "card_token", "required")
		}
		if in.PaymentMethodID == nil || *in.PaymentMethodID == "" {
			return fielderr.On(errs.ErrPaymentMethodIDRequired, "payment_method_id", "required")
		}
		if in.Installments != nil && *in.Installments < 1 {
			return fielderr.On(errs.ErrInstallmentsBelowOne, "installments", "min=1")
		}
		return nil
	default:
		return fielderr.On(errs.ErrInvalidPaymentMethod, "payment_method", "one of pix, credit_card, debit_card")
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
		return nil, fielderr.On(errs.ErrNoItems, "items", "min=1")
	}

	seenProduct := make(map[uuid.UUID]bool, len(items))
	seenProgram := make(map[uuid.UUID]bool, len(items))
	seenAttendee := make(map[string]bool, len(items))
	hasTicket := false

	for i, line := range items {
		field := fmt.Sprintf("items[%d]", i)
		if line.Quantity < 1 {
			return nil, fielderr.On(errs.ErrQuantityBelowOne, field+".quantity", "min=1")
		}
		switch line.ItemType {
		case models.PurchaseItemTypeTicket:
			if line.Quantity != 1 {
				return nil, fielderr.On(errs.ErrTicketQuantity, field+".quantity", "tickets are one per person — send one line per attendee")
			}
			if line.Attendee == nil {
				return nil, fielderr.On(errs.ErrTicketNeedsAttendee, field+".attendee", "required")
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
				return nil, fielderr.On(errs.ErrOneTicketPerPerson, field+".attendee", "this person already has a ticket in the cart")
			}
			seenAttendee[attendeeKey] = true
			hasTicket = true
		case models.PurchaseItemTypeProduct:
			if seenProduct[line.ItemID] {
				return nil, fielderr.On(errs.ErrDuplicateProductLine, field, "this product is already in the cart")
			}
			seenProduct[line.ItemID] = true
		case models.PurchaseItemTypeProgramOccurrence:
			if line.Quantity != 1 {
				return nil, fielderr.On(errs.ErrProgramQuantity, field+".quantity", "program lines are quantity 1")
			}
			if seenProgram[line.ItemID] {
				return nil, fielderr.On(errs.ErrDuplicateProgramLine, field, "this program occurrence is already in the cart")
			}
			seenProgram[line.ItemID] = true
		default:
			return nil, fielderr.On(errs.ErrInvalidItemType, field+".item_type", "one of ticket, product, program_occurrence")
		}
	}

	if len(seenProgram) > 0 && !hasTicket {
		return nil, fielderr.On(errs.ErrProgramNeedsTicketInCart, "items", "a participation attaches to a ticket's registration")
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
			return nil, fielderr.On(errs.ErrTicketNeedsAttendeeEmail, fmt.Sprintf("items[%d].attendee.email", i), "required")
		}
		line.Attendee.Email = email

		// One lookup per attendee: ErrActorNotFound means the email has no
		// account in the univents identityx project.
		actor, err := o.Actors.GetByEmail(ctx, email)
		switch {
		case err != nil && !errors.Is(err, ports.ErrActorNotFound):
			// IdentityX unreachable — fail the checkout closed: an attendee
			// whose account cannot be verified must not be silently gifted
			// as email-only (they may already hold a ticket).
			return nil, err
		case errors.Is(err, ports.ErrActorNotFound) && line.Attendee.UserID != nil:
			// A user_id was claimed for an email with no account — the pair
			// cannot both be true.
			return nil, fielderr.On(errs.ErrAttendeeEmailUnknown, fmt.Sprintf("items[%d].attendee.email", i), email)
		case errors.Is(err, ports.ErrActorNotFound):
			// No account yet — email-only gift; keep the nil user id.
		case line.Attendee.UserID == nil:
			line.Attendee.UserID = &actor.ID
		case actor.ID != *line.Attendee.UserID:
			return nil, fielderr.On(errs.ErrAttendeeUserMismatch, fmt.Sprintf("items[%d].attendee.user_id", i), line.Attendee.UserID.String())
		}

		// Post-resolution dedup: two lines that resolve to the same person
		// (explicit id, or a resolved account) collide here; email-only
		// lines collide on the normalized email.
		key := "e:" + line.Attendee.Email
		if line.Attendee.UserID != nil {
			key = "u:" + line.Attendee.UserID.String()
		}
		if seen[key] {
			return nil, fielderr.On(errs.ErrOneTicketPerPerson, fmt.Sprintf("items[%d].attendee", i), "this person already has a ticket in the cart")
		}
		seen[key] = true
	}
	return out, nil
}
