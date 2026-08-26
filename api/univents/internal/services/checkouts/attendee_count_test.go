package checkouts_test

import (
	"context"
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	"univents/models"
)

// seedRegistration inserts one registration row for the edition (via the
// repo's Create — no availability checks, the point is the count read, not
// the buy flow).
func seedRegistration(t *testing.T, r regRepo, editionID, ticketID uuid.UUID, status models.RegistrationStatus) {
	t.Helper()
	// One attendee per email: the active-per-edition email index enforces
	// one ticket per person, and each seeded row is a different person.
	attendeeID := uuid.New()
	_, err := r.Create(context.Background(), &models.Registration{
		EditionID:      editionID,
		TicketTypeID:   ticketID,
		PurchaserID:    attendeeID,
		AttendeeUserID: &attendeeID,
		AttendeeEmail:  "attendee-" + attendeeID.String() + "@example.com",
		AttendeeName:   "Attendee",
		Status:         status,
	})
	if err != nil {
		t.Fatalf("seed registration (%s): %v", status, err)
	}
}

// regRepo is the slice of the registration repo the seeding helper needs —
// satisfied by *repos.Repos' Registrations field.
type regRepo interface {
	Create(ctx context.Context, toCreate *models.Registration) (*models.Registration, error)
}

// TestAttendeeCount_CountsConfirmedOnly pins the count semantics: only
// confirmed (paid) registrations count; pending reservations and
// cancelled/expired rows do not.
func TestAttendeeCount_CountsConfirmedOnly(t *testing.T) {
	r, ops, _, _ := newOps(t, nil)
	fx := seedStore(t, r)

	seedRegistration(t, r.Registrations, fx.editionID, fx.ticketID, models.RegistrationStatusConfirmed)
	seedRegistration(t, r.Registrations, fx.editionID, fx.ticketID, models.RegistrationStatusConfirmed)
	seedRegistration(t, r.Registrations, fx.editionID, fx.ticketID, models.RegistrationStatusConfirmed)
	seedRegistration(t, r.Registrations, fx.editionID, fx.ticketID, models.RegistrationStatusPending)
	seedRegistration(t, r.Registrations, fx.editionID, fx.ticketID, models.RegistrationStatusCancelled)
	seedRegistration(t, r.Registrations, fx.editionID, fx.ticketID, models.RegistrationStatusExpired)

	count, err := ops.AttendeeCount(context.Background(), fx.editionID)
	if err != nil {
		t.Fatalf("AttendeeCount: %v", err)
	}
	if count != 3 {
		t.Fatalf("count = %d, want 3 (confirmed only)", count)
	}
}

// TestAttendeeCount_OtherEditionExcluded pins edition scoping: registrations
// of a sibling edition never leak into the count.
func TestAttendeeCount_OtherEditionExcluded(t *testing.T) {
	r, ops, _, _ := newOps(t, nil)
	fx := seedStore(t, r)
	other := seedStore(t, r)

	seedRegistration(t, r.Registrations, fx.editionID, fx.ticketID, models.RegistrationStatusConfirmed)
	seedRegistration(t, r.Registrations, other.editionID, other.ticketID, models.RegistrationStatusConfirmed)
	seedRegistration(t, r.Registrations, other.editionID, other.ticketID, models.RegistrationStatusConfirmed)

	count, err := ops.AttendeeCount(context.Background(), fx.editionID)
	if err != nil {
		t.Fatalf("AttendeeCount: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1 (sibling edition excluded)", count)
	}
}

// TestAttendeeCount_EmptyEditionIsZero pins that an edition with no
// confirmed registrations is 0, not an error.
func TestAttendeeCount_EmptyEditionIsZero(t *testing.T) {
	r, ops, _, _ := newOps(t, nil)
	fx := seedStore(t, r)

	count, err := ops.AttendeeCount(context.Background(), fx.editionID)
	if err != nil {
		t.Fatalf("AttendeeCount: %v", err)
	}
	if count != 0 {
		t.Fatalf("count = %d, want 0", count)
	}
}

// TestAttendeeCount_UnknownEditionNotFound pins the 404 contract: the count
// read exists, the edition does not.
func TestAttendeeCount_UnknownEditionNotFound(t *testing.T) {
	r, ops, _, _ := newOps(t, nil)
	seedStore(t, r)

	_, err := ops.AttendeeCount(context.Background(), uuid.New())
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("err = %v, want NOT_FOUND", err)
	}
}
