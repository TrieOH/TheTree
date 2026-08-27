package checkouts

import (
	"context"
	"errors"
	"strings"

	"lib/telemetry"
	"univents/models"
	"univents/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// MyTicket returns the caller's active ticket for the edition — their own
// registration where attendee_user_id = caller and status is pending or
// confirmed — with its ticket type. Returns (nil, nil) when the caller
// holds none (cancelled/expired registrations do not count). Read-only;
// never transitions state. An unknown edition is NOT_FOUND (the endpoint
// exists, the edition does not).
//
// This is also the lazy gift-claim surface: a caller who holds no ticket
// under their id may be a gifted recipient who just created their account
// — if their account email matches an email-only gift in this edition, the
// registration is tied to their id (the claim) and, when confirmed, the
// deferred badge is emitted (the gift email told them to create the account
// under that email, so the email is the proof).
func (o *Operations) MyTicket(ctx context.Context, editionID, userID uuid.UUID) (*models.MyTicket, error) {
	ctx, span := telemetry.StartSpan(ctx, "CheckoutsService.MyTicket")
	defer span.End()

	_, err := o.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	reg, err := o.activeRegistrationFor(ctx, editionID, userID)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		reg, err = o.claimGiftByAccountEmail(ctx, editionID, userID)
		if err != nil {
			return nil, err
		}
		if reg == nil {
			//nolint:nilnil // no ticket is a normal state (front falls back to the buy flow)
			return nil, nil
		}
	}

	ticketType, err := o.ticketTypes.GetByID(ctx, reg.TicketTypeID)
	if err != nil {
		return nil, err
	}
	return &models.MyTicket{
		RegistrationID: reg.ID,
		TicketType:     *ticketType,
		Status:         reg.Status,
	}, nil
}

// claimGiftByAccountEmail claims the caller's email-only gifted
// registration in the edition, if any: the caller's account email (from
// IdentityX) is matched against active gifts (attendee_user_id NULL) and
// the registration is tied to their actor id. A confirmed claim emits the
// deferred badge. Idempotent: once claimed, the id is set and the my-ticket
// read finds the registration directly (this path no longer runs). No gift
// → nil, nil.
func (o *Operations) claimGiftByAccountEmail(ctx context.Context, editionID, userID uuid.UUID) (*models.Registration, error) {
	actor, err := o.actors.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ports.ErrActorNotFound) {
			//nolint:nilnil // unknown actor — nothing to claim by
			return nil, nil
		}
		return nil, err
	}
	if actor.Email == nil {
		//nolint:nilnil // no email on the account — cannot match a gift
		return nil, nil
	}
	email := strings.TrimSpace(strings.ToLower(*actor.Email))
	if email == "" {
		//nolint:nilnil // empty email — cannot match a gift
		return nil, nil
	}

	claimed, err := o.registrations.ClaimByAttendeeEmail(ctx, editionID, email, userID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			//nolint:nilnil // no active gift for this email — a normal state
			return nil, nil
		}
		return nil, err
	}

	if claimed.Status == models.RegistrationStatusConfirmed {
		// The gift is paid — emit the badge that was deferred at approval
		// (EmitForConfirmedRegistration skips registrations with no
		// account; now that the id is tied, it emits + emails).
		_, err := o.badges.EmitForConfirmedRegistration(ctx, claimed.ID)
		if err != nil {
			return nil, err
		}
	}
	return claimed, nil
}
