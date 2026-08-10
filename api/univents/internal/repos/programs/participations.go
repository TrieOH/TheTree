package programs

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

// CreateParticipation registers a program participation at checkout
// (split 7), attached to the required ticket item's registration — one of
// the purchase's materialized rows (D4).
func (repo *Repo) CreateParticipation(ctx context.Context, toCreate *models.ProgramParticipation) (*models.ProgramParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.CreateParticipation")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreateProgramParticipation(ctx, sqlc.CreateProgramParticipationParams{
		EditionID:      toCreate.EditionID,
		OccurrenceID:   toCreate.OccurrenceID,
		RegistrationID: toCreate.RegistrationID,
		Status:         string(toCreate.Status),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapParticipation(result)), nil
}

// UpdateParticipationStatus flips a participation (registered→cancelled/
// attended/no_show) — the webhook receiver (split 4) and the expiry worker
// (split 7) write side.
func (repo *Repo) UpdateParticipationStatus(ctx context.Context, id uuid.UUID, status models.ProgramParticipationStatus) (*models.ProgramParticipation, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.UpdateParticipationStatus")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).UpdateProgramParticipationStatus(ctx, sqlc.UpdateProgramParticipationStatusParams{
		ID:     id,
		Status: string(status),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapParticipation(result)), nil
}
