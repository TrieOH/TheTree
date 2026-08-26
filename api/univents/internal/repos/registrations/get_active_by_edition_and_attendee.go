package registrations

import (
	"context"

	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

// GetActiveByEditionAndAttendee returns the attendee's active (pending or
// confirmed) registration in an edition — the one-ticket-per-person check
// (checkout) and the my-ticket read. An absent registration surfaces as
// NOT_FOUND (the shared ErrorHandler maps ErrNoRows); the service layer
// interprets "not found" as "holds no ticket", which is a normal state
// there, not an error here.
func (repo *Repo) GetActiveByEditionAndAttendee(ctx context.Context, editionID, attendeeID uuid.UUID) (*models.Registration, error) {
	ctx, span := telemetry.StartSpan(ctx, "RegistrationsRepo.GetActiveByEditionAndAttendee")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetActiveByEditionAndAttendee(ctx, sqlc.GetActiveByEditionAndAttendeeParams{
		EditionID:      editionID,
		AttendeeUserID: &attendeeID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapRegistration(row)), nil
}
