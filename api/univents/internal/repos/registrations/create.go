package registrations

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

// Create inserts a pending registration at checkout (split 7) — one of the
// purchase's materialized rows (D4).
func (repo *Repo) Create(ctx context.Context, toCreate *models.Registration) (*models.Registration, error) {
	ctx, span := telemetry.StartSpan(ctx, "RegistrationsRepo.Create")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreateRegistration(ctx, sqlc.CreateRegistrationParams{
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
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapRegistration(result)), nil
}

// UpdateStatus flips a registration (pending→confirmed/cancelled/expired) —
// the webhook receiver (split 4) and the expiry worker (split 7) write side.
func (repo *Repo) UpdateStatus(ctx context.Context, id uuid.UUID, status models.RegistrationStatus, reason *string) (*models.Registration, error) {
	ctx, span := telemetry.StartSpan(ctx, "RegistrationsRepo.UpdateStatus")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdateRegistrationStatus(ctx, sqlc.UpdateRegistrationStatusParams{
		ID:           id,
		Status:       string(status),
		StatusReason: reason,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapRegistration(result)), nil
}
