package registrations

import (
	"context"

	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

// GetActiveByEditionAndAttendeeEmail returns the email-only recipient's
// active (pending or confirmed) registration in an edition — the
// one-ticket-per-person check for gifted tickets to recipients without an
// IdentityX account (attendee_user_id NULL, checkout). An absent
// registration surfaces as NOT_FOUND (the shared ErrorHandler maps
// ErrNoRows); the service layer interprets "not found" as "holds no
// ticket", which is a normal state there, not an error here.
func (repo *Repo) GetActiveByEditionAndAttendeeEmail(ctx context.Context, editionID uuid.UUID, attendeeEmail string) (*models.Registration, error) {
	ctx, span := telemetry.StartSpan(ctx, "RegistrationsRepo.GetActiveByEditionAndAttendeeEmail")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetActiveByEditionAndAttendeeEmail(ctx, sqlc.GetActiveByEditionAndAttendeeEmailParams{
		EditionID:     editionID,
		AttendeeEmail: attendeeEmail,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapRegistration(row)), nil
}
