package checkouts_test

import (
	"context"
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	"univents/internal/services/checkouts"
	"univents/models"

	"sdk/payssage"
)

// seedFreeTicket adds a zero-price ticket type so a checkout confirms
// immediately (approved + confirmed registration) without an intent.
func seedFreeTicket(t *testing.T, r interface {
	Create(ctx context.Context, toCreate *models.TicketType) (*models.TicketType, error)
}, editionID uuid.UUID) *models.TicketType {
	t.Helper()
	ticket, err := r.Create(context.Background(), &models.TicketType{
		EditionID:   editionID,
		Name:        "Free",
		AccessLevel: 0,
		PriceCents:  0,
		MaxQuantity: new(int(1)),
	})
	if err != nil {
		t.Fatalf("seed free ticket: %v", err)
	}
	return ticket
}

// TestMyTicket_ReturnsHeldTicket pins the read: a confirmed registration
// comes back with its ticket type so the front can compare prices.
func TestMyTicket_ReturnsHeldTicket(t *testing.T) {
	r, ops, fs, _ := newOps(t, func(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return pixIntent(), nil
	})
	fx := seedStore(t, r)
	buyerID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	freeTicket := seedFreeTicket(t, r.TicketTypes, fx.editionID)
	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID,
		pixInput(ticketLine(freeTicket.ID, selfAttendee(buyerID))))
	if err != nil {
		t.Fatalf("free checkout: %v", err)
	}

	ticket, err := ops.MyTicket(context.Background(), fx.editionID, buyerID)
	if err != nil {
		t.Fatalf("MyTicket: %v", err)
	}
	if ticket == nil {
		t.Fatal("ticket = nil, want the held ticket")
	}
	if ticket.Status != models.RegistrationStatusConfirmed {
		t.Fatalf("status = %s, want confirmed", ticket.Status)
	}
	if ticket.TicketType.ID != freeTicket.ID {
		t.Fatalf("ticket_type id = %s, want %s", ticket.TicketType.ID, freeTicket.ID)
	}
	if ticket.TicketType.PriceCents != 0 {
		t.Fatalf("ticket_type price = %d, want 0", ticket.TicketType.PriceCents)
	}
	if ticket.TicketType.Name != "Free" {
		t.Fatalf("ticket_type name = %s, want Free", ticket.TicketType.Name)
	}
}

// TestMyTicket_PendingCountsAsHeld pins that an unpaid reservation is still
// a held ticket (the front must not offer a fresh buy while one is pending).
func TestMyTicket_PendingCountsAsHeld(t *testing.T) {
	r, ops, fs, _ := newOps(t, func(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return pixIntent(), nil
	})
	fx := seedStore(t, r)
	buyerID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID,
		pixInput(ticketLine(fx.ticketID, selfAttendee(buyerID))))
	if err != nil {
		t.Fatalf("checkout: %v", err)
	}

	ticket, err := ops.MyTicket(context.Background(), fx.editionID, buyerID)
	if err != nil {
		t.Fatalf("MyTicket: %v", err)
	}
	if ticket == nil {
		t.Fatal("ticket = nil, want the pending reservation")
	}
	if ticket.Status != models.RegistrationStatusPending {
		t.Fatalf("status = %s, want pending", ticket.Status)
	}
}

// TestMyTicket_GiftedTicketCounts pins that a ticket gifted to the caller
// (attendee != purchaser) is still their ticket — one per person, whoever
// paid.
func TestMyTicket_GiftedTicketCounts(t *testing.T) {
	r, ops, fs, _ := newOps(t, func(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return pixIntent(), nil
	})
	fx := seedStore(t, r)
	buyerID := uuid.New()
	recipientID := uuid.New()
	fs.actors.seed("friend@example.com", recipientID)

	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{UserID: &recipientID, Email: "friend@example.com", Name: "John Doe"}),
	))
	if err != nil {
		t.Fatalf("checkout: %v", err)
	}

	// The recipient holds the ticket; the purchaser does not.
	ticket, err := ops.MyTicket(context.Background(), fx.editionID, recipientID)
	if err != nil {
		t.Fatalf("MyTicket(recipient): %v", err)
	}
	if ticket == nil {
		t.Fatal("recipient ticket = nil, want the gifted ticket")
	}
	if ticket.TicketType.ID != fx.ticketID {
		t.Fatalf("ticket_type id = %s, want %s", ticket.TicketType.ID, fx.ticketID)
	}

	purchaserTicket, err := ops.MyTicket(context.Background(), fx.editionID, buyerID)
	if err != nil {
		t.Fatalf("MyTicket(purchaser): %v", err)
	}
	if purchaserTicket != nil {
		t.Fatalf("purchaser ticket = %+v, want nil (they gifted, they don't hold)", purchaserTicket)
	}
}

// TestMyTicket_NoTicketIsNil pins the no-ticket contract: a caller who
// never bought returns nil (the front falls back to the normal buy flow).
func TestMyTicket_NoTicketIsNil(t *testing.T) {
	r, ops, _, _ := newOps(t, nil)
	fx := seedStore(t, r)

	ticket, err := ops.MyTicket(context.Background(), fx.editionID, uuid.New())
	if err != nil {
		t.Fatalf("MyTicket: %v", err)
	}
	if ticket != nil {
		t.Fatalf("ticket = %+v, want nil", ticket)
	}
}

// TestMyTicket_UnknownEditionNotFound pins the 404: a garbage edition id is
// not "no ticket", it is an unknown resource.
func TestMyTicket_UnknownEditionNotFound(t *testing.T) {
	r, ops, _, _ := newOps(t, nil)
	_ = r
	_, err := ops.MyTicket(context.Background(), uuid.New(), uuid.New())
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("err = %v, want NOT_FOUND", err)
	}
}

// TestMyTicket_ClaimsGiftByEmail pins the lazy gift claim: an accountless
// recipient who was gifted a (free, confirmed) ticket and then created an
// account — their email matches the email-only gift, so the my-ticket read
// ties the registration to their id and emits the badge that was deferred
// at confirmation.
func TestMyTicket_ClaimsGiftByEmail(t *testing.T) {
	r, ops, fs, _ := newOps(t, func(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return pixIntent(), nil
	})
	fx := seedStore(t, r)
	buyerID := uuid.New()
	recipientID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	freeTicket := seedFreeTicket(t, r.TicketTypes, fx.editionID)
	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(freeTicket.ID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
	))
	if err != nil {
		t.Fatalf("gift checkout: %v", err)
	}
	// The gift email job was enqueued at confirmation; the badge was not
	// emitted for real (the free-order emit call happens, but the real
	// service skips accountless registrations).
	if len(fs.river.giftEmailArgs()) != 1 {
		t.Fatalf("gifts.send_email jobs = %d, want 1", len(fs.river.giftEmailArgs()))
	}

	// The recipient created their account under the gifted email after the
	// gift — the claim matches by the account's own email.
	fs.actors.seed("friend@example.com", recipientID)

	// The recipient created their account and loads the edition: my-ticket
	// claims the gift, emits the deferred badge, and returns the ticket.
	ticket, err := ops.MyTicket(context.Background(), fx.editionID, recipientID)
	if err != nil {
		t.Fatalf("MyTicket(recipient): %v", err)
	}
	if ticket == nil || ticket.TicketType.ID != freeTicket.ID {
		t.Fatalf("ticket = %+v, want the gifted free ticket", ticket)
	}
	if ticket.Status != models.RegistrationStatusConfirmed {
		t.Fatalf("status = %s, want confirmed", ticket.Status)
	}

	reg, err := r.Registrations.GetByID(context.Background(), ticket.RegistrationID)
	if err != nil {
		t.Fatalf("load registration: %v", err)
	}
	if reg.AttendeeUserID == nil || *reg.AttendeeUserID != recipientID {
		t.Fatalf("attendee_user_id = %v, want the recipient %s (claimed)", reg.AttendeeUserID, recipientID)
	}

	// The claim emitted the deferred badge: 1 call at checkout (skipped by
	// the real service for accountless) + 1 at claim.
	if got := fs.badges.emittedCount(); got != 2 {
		t.Fatalf("badge emits = %d, want 2 (checkout call + claim emit)", got)
	}
}

// TestMyTicket_ClaimIsIdempotent pins that a second my-ticket read does not
// re-claim or re-emit: once the registration carries the id, the read finds
// it directly and the claim path never runs again.
func TestMyTicket_ClaimIsIdempotent(t *testing.T) {
	r, ops, fs, _ := newOps(t, func(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return pixIntent(), nil
	})
	fx := seedStore(t, r)
	buyerID := uuid.New()
	recipientID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	freeTicket := seedFreeTicket(t, r.TicketTypes, fx.editionID)
	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(freeTicket.ID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
	))
	if err != nil {
		t.Fatalf("gift checkout: %v", err)
	}
	fs.actors.seed("friend@example.com", recipientID)

	_, err = ops.MyTicket(context.Background(), fx.editionID, recipientID)
	if err != nil {
		t.Fatalf("MyTicket (claim): %v", err)
	}
	before := fs.badges.emittedCount()
	_, err = ops.MyTicket(context.Background(), fx.editionID, recipientID)
	if err != nil {
		t.Fatalf("MyTicket (repeat): %v", err)
	}
	if after := fs.badges.emittedCount(); after != before {
		t.Fatalf("badge emits = %d → %d, want no re-emit on repeat reads", before, after)
	}
}

// TestMyTicket_PendingGiftClaimsWithoutBadge pins that claiming an unpaid
// (pending) gift ties the id but does not emit — the badge arrives when the
// webhook confirms the payment (EmitForConfirmedRegistration handles the
// confirmed registration normally once the id is set).
func TestMyTicket_PendingGiftClaimsWithoutBadge(t *testing.T) {
	r, ops, fs, _ := newOps(t, func(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return pixIntent(), nil
	})
	fx := seedStore(t, r)
	buyerID := uuid.New()
	recipientID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	// Paid gift: stays pending (only the webhook confirms).
	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
	))
	if err != nil {
		t.Fatalf("gift checkout: %v", err)
	}
	fs.actors.seed("friend@example.com", recipientID)

	ticket, err := ops.MyTicket(context.Background(), fx.editionID, recipientID)
	if err != nil {
		t.Fatalf("MyTicket(recipient): %v", err)
	}
	if ticket == nil {
		t.Fatal("ticket = nil, want the claimed pending gift")
	}
	if ticket.Status != models.RegistrationStatusPending {
		t.Fatalf("status = %s, want pending", ticket.Status)
	}
	if got := fs.badges.emittedCount(); got != 0 {
		t.Fatalf("badge emits = %d, want 0 (payment not confirmed yet)", got)
	}
}

// TestMyTicket_NoGiftIsNotClaimed pins that a ticketless caller with an
// account but no matching gift gets nil, and no registration is created or
// tied.
func TestMyTicket_NoGiftIsNotClaimed(t *testing.T) {
	r, ops, fs, _ := newOps(t, nil)
	fx := seedStore(t, r)
	bystanderID := uuid.New()
	fs.actors.seed("bystander@example.com", bystanderID)

	ticket, err := ops.MyTicket(context.Background(), fx.editionID, bystanderID)
	if err != nil {
		t.Fatalf("MyTicket: %v", err)
	}
	if ticket != nil {
		t.Fatalf("ticket = %+v, want nil (no gift to claim)", ticket)
	}
	if got := fs.badges.emittedCount(); got != 0 {
		t.Fatalf("badge emits = %d, want 0", got)
	}
}
