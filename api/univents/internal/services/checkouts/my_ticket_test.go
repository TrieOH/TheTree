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
	r, ops, _, _ := newOps(t, func(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return pixIntent(), nil
	})
	fx := seedStore(t, r)
	buyerID := uuid.New()

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
	r, ops, _, _ := newOps(t, func(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return pixIntent(), nil
	})
	fx := seedStore(t, r)
	buyerID := uuid.New()

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
	r, ops, _, _ := newOps(t, func(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return pixIntent(), nil
	})
	fx := seedStore(t, r)
	buyerID := uuid.New()
	recipientID := uuid.New()

	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{UserID: recipientID, Email: "friend@example.com", Name: "John Doe"}),
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
