package repo

import (
	"context"
	"lib/database"
	"univents/errs"
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.RegistrationRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrRegistrationNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapRegistration(src sqlc.Registration) models.Registration {
	return models.Registration{
		ID:               src.ID,
		EditionID:        src.EditionID,
		TicketTypeID:     src.TicketTypeID,
		PurchaserID:      src.PurchaserID,
		AttendeeUserID:   src.AttendeeUserID,
		AttendeeEmail:    src.AttendeeEmail,
		AttendeeName:     src.AttendeeName,
		Status:           models.RegistrationStatus(src.Status),
		StatusReason:     src.StatusReason,
		PayssageIntentID: src.PayssageIntentID,
		CreatedAt:        src.CreatedAt,
		UpdatedAt:        src.UpdatedAt,
		DeletedAt:        src.DeletedAt,
	}
}

// ClaimAllByAttendeeEmail ties every active email-only gift for the email
// to the recipient's IdentityX account — the profile-badges claim, fired
// from the badges read when a user's profile loads (a profile requires an
// account, so the account email deterministically resolves the gifts).
// Crosses editions, unlike ClaimByAttendeeEmail. Returns the claimed
// registrations so the caller can emit deferred badges for the confirmed
// ones. No gifts → empty slice, nil.
func (repo *Repo) ClaimAllByAttendeeEmail(ctx context.Context, attendeeEmail string, userID uuid.UUID) ([]*models.Registration, error) {
	rows, err := repo.db(ctx).ClaimAllByAttendeeEmail(ctx, sqlc.ClaimAllByAttendeeEmailParams{
		AttendeeEmail:  attendeeEmail,
		AttendeeUserID: &userID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	claimed := make([]*models.Registration, 0, len(rows))
	for _, row := range rows {
		claimed = append(claimed, new(mapRegistration(row)))
	}
	return claimed, nil
}

// ClaimByAttendeeEmail ties the email-only gifted registration (attendee_user_id
// NULL) to the recipient's IdentityX account — the gift claim, fired lazily
// from the my-ticket read when the caller holds no ticket under their own id
// but their account email matches a gifted one. Only active rows are
// claimed; an absent match is ErrRegistrationNotFound
// (database.ErrNotFound), which the service treats as "no gift to claim".
func (repo *Repo) ClaimByAttendeeEmail(ctx context.Context, editionID uuid.UUID, attendeeEmail string, userID uuid.UUID) (*models.Registration, error) {
	row, err := repo.db(ctx).ClaimRegistrationByEmail(ctx, sqlc.ClaimRegistrationByEmailParams{
		EditionID:      editionID,
		AttendeeEmail:  attendeeEmail,
		AttendeeUserID: &userID,
	})
	return database.One(row, err, repo.dbe, mapRegistration)
}

// CountConfirmedByEdition returns the number of confirmed attendees of an
// edition — the attendee-count read (GET /editions/{id}/attendees/count).
// Only paid registrations count; pending reservations and cancelled/expired
// rows do not. An edition with no confirmed registrations is 0, not an
// error.
func (repo *Repo) CountConfirmedByEdition(ctx context.Context, editionID uuid.UUID) (int64, error) {
	count, err := repo.db(ctx).CountConfirmedByEdition(ctx, editionID)
	if err != nil {
		return 0, repo.dbe(err)
	}
	return count, nil
}

// Create inserts a pending registration at checkout (split 7) — one of the
// purchase's materialized rows (D4).
func (repo *Repo) Create(ctx context.Context, toCreate *models.Registration) (*models.Registration, error) {
	result, err := repo.db(ctx).CreateRegistration(ctx, sqlc.CreateRegistrationParams{
		EditionID:        toCreate.EditionID,
		TicketTypeID:     toCreate.TicketTypeID,
		PurchaserID:      toCreate.PurchaserID,
		AttendeeUserID:   toCreate.AttendeeUserID,
		AttendeeEmail:    toCreate.AttendeeEmail,
		AttendeeName:     toCreate.AttendeeName,
		Status:           string(toCreate.Status),
		StatusReason:     toCreate.StatusReason,
		PayssageIntentID: toCreate.PayssageIntentID,
	})
	return database.One(result, err, repo.dbe, mapRegistration)
}

// UpdateStatus flips a registration (pending→confirmed/cancelled/expired) —
// the webhook receiver (split 4) and the expiry worker (split 7) write side.
func (repo *Repo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.RegistrationStatus, reason *string) (*models.Registration, error) {
	result, err := repo.db(ctx).UpdateRegistrationStatus(ctx, sqlc.UpdateRegistrationStatusParams{
		ID:           id,
		Status:       string(status),
		StatusReason: reason,
	})
	return database.One(result, err, repo.dbe, mapRegistration)
}

// GetActiveByEditionAndAttendee returns the attendee's active (pending or
// confirmed) registration in an edition — the one-ticket-per-person check
// (checkout) and the my-ticket read. An absent registration is
// ErrRegistrationNotFound (database.ErrNotFound); the service layer
// interprets "not found" as "holds no ticket", which is a normal state
// there, not an error here.
func (repo *Repo) GetActiveByEditionAndAttendee(ctx context.Context, editionID, attendeeID uuid.UUID) (*models.Registration, error) {
	row, err := repo.db(ctx).GetActiveByEditionAndAttendee(ctx, sqlc.GetActiveByEditionAndAttendeeParams{
		EditionID:      editionID,
		AttendeeUserID: &attendeeID,
	})
	return database.One(row, err, repo.dbe, mapRegistration)
}

// GetActiveByEditionAndAttendeeEmail returns the email-only recipient's
// active (pending or confirmed) registration in an edition — the
// one-ticket-per-person check for gifted tickets to recipients without an
// IdentityX account (attendee_user_id NULL, checkout). An absent
// registration is ErrRegistrationNotFound (database.ErrNotFound); the
// service layer interprets "not found" as "holds no
// ticket", which is a normal state there, not an error here.
func (repo *Repo) GetActiveByEditionAndAttendeeEmail(ctx context.Context, editionID uuid.UUID, attendeeEmail string) (*models.Registration, error) {
	row, err := repo.db(ctx).GetActiveByEditionAndAttendeeEmail(ctx, sqlc.GetActiveByEditionAndAttendeeEmailParams{
		EditionID:     editionID,
		AttendeeEmail: attendeeEmail,
	})
	return database.One(row, err, repo.dbe, mapRegistration)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Registration, error) {
	row, err := repo.db(ctx).GetRegistrationByID(ctx, id)
	return database.One(row, err, repo.dbe, mapRegistration)
}
