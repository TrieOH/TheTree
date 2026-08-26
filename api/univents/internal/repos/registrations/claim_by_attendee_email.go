package registrations

import (
	"context"

	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

// ClaimByAttendeeEmail ties the email-only gifted registration (attendee_user_id
// NULL) to the recipient's IdentityX account — the gift claim, fired lazily
// from the my-ticket read when the caller holds no ticket under their own id
// but their account email matches a gifted one. Only active rows are
// claimed; an absent match surfaces as NOT_FOUND (the shared ErrorHandler
// maps ErrNoRows), which the service treats as "no gift to claim".
func (repo *Repo) ClaimByAttendeeEmail(ctx context.Context, editionID uuid.UUID, attendeeEmail string, userID uuid.UUID) (*models.Registration, error) {
	ctx, span := telemetry.StartSpan(ctx, "RegistrationsRepo.ClaimByAttendeeEmail")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).ClaimRegistrationByEmail(ctx, sqlc.ClaimRegistrationByEmailParams{
		EditionID:      editionID,
		AttendeeEmail:  attendeeEmail,
		AttendeeUserID: &userID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapRegistration(row)), nil
}
