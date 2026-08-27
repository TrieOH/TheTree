package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

// RegistrationRepo is the registration read surface. The checkout feature owns
// the write side; readers (badges, certifications) consume confirmed rows.
type RegistrationRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Registration, error)

	// CountConfirmedByEdition returns the number of confirmed attendees of
	// an edition — the attendee-count read. Pending reservations and
	// cancelled/expired rows do not count.
	CountConfirmedByEdition(ctx context.Context, editionID uuid.UUID) (int64, error)

	// GetActiveByEditionAndAttendee returns the attendee's active (pending or
	// confirmed) registration in an edition, or nil when they hold none — the
	// one-ticket-per-person check (checkout) and the my-ticket read.
	GetActiveByEditionAndAttendee(ctx context.Context, editionID, attendeeID uuid.UUID) (*models.Registration, error)

	// GetActiveByEditionAndAttendeeEmail returns the email-only recipient's
	// active (pending or confirmed) registration in an edition, or nil when
	// they hold none — the one-ticket-per-person check for gifted tickets to
	// recipients without an IdentityX account (attendee_user_id NULL).
	GetActiveByEditionAndAttendeeEmail(ctx context.Context, editionID uuid.UUID, attendeeEmail string) (*models.Registration, error)

	// ClaimByAttendeeEmail ties an email-only gifted registration to the
	// recipient's IdentityX account (gift claim, fired lazily from the
	// my-ticket read). NOT_FOUND when the email holds no active gift.
	ClaimByAttendeeEmail(ctx context.Context, editionID uuid.UUID, attendeeEmail string, userID uuid.UUID) (*models.Registration, error)

	// ClaimAllByAttendeeEmail ties every active email-only gift for the
	// email to the recipient's IdentityX account (the profile-badges claim,
	// fired from the badges read) and returns the claimed registrations so
	// the caller can emit deferred badges for the confirmed ones. No gifts →
	// empty slice, nil.
	ClaimAllByAttendeeEmail(ctx context.Context, attendeeEmail string, userID uuid.UUID) ([]*models.Registration, error)

	// Create + UpdateStatus are the checkout (split 7) / webhook receiver
	// (split 4) write side: pending rows are materialized at checkout and
	// flipped on approve/cancel/expire.
	Create(ctx context.Context, toCreate *models.Registration) (*models.Registration, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.RegistrationStatus, reason *string) (*models.Registration, error)
}
