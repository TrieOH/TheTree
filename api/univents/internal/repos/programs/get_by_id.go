package programs

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.GetByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetProgramByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProgram(result)), nil
}

func (repo *Repo) GetOccurrenceByID(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.GetOccurrenceByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetProgramOccurrenceByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProgramOccurrence(result)), nil
}
