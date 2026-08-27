package registrations

import (
	"context"

	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

// ClaimAllByAttendeeEmail ties every active email-only gift for the email
// to the recipient's IdentityX account — the profile-badges claim, fired
// from the badges read when a user's profile loads (a profile requires an
// account, so the account email deterministically resolves the gifts).
// Crosses editions, unlike ClaimByAttendeeEmail. Returns the claimed
// registrations so the caller can emit deferred badges for the confirmed
// ones. No gifts → empty slice, nil.
func (repo *Repo) ClaimAllByAttendeeEmail(ctx context.Context, attendeeEmail string, userID uuid.UUID) ([]*models.Registration, error) {
	ctx, span := telemetry.StartSpan(ctx, "RegistrationsRepo.ClaimAllByAttendeeEmail")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ClaimAllByAttendeeEmail(ctx, sqlc.ClaimAllByAttendeeEmailParams{
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
