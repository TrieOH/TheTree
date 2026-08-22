package registrations

import (
	"context"

	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

// CountConfirmedByEdition returns the number of confirmed attendees of an
// edition — the attendee-count read (GET /editions/{id}/attendees/count).
// Only paid registrations count; pending reservations and cancelled/expired
// rows do not. An edition with no confirmed registrations is 0, not an
// error.
func (repo *Repo) CountConfirmedByEdition(ctx context.Context, editionID uuid.UUID) (int64, error) {
	ctx, span := telemetry.StartSpan(ctx, "RegistrationsRepo.CountConfirmedByEdition")
	defer span.End()
	count, err := database.Queries(ctx, repo.q).CountConfirmedByEdition(ctx, editionID)
	if err != nil {
		return 0, repo.dbe(err)
	}
	return count, nil
}
