package programs

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.Delete")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).DeleteProgram(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProgram(result)), nil
}

func (repo *Repo) DeleteOccurrence(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.DeleteOccurrence")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).DeleteProgramOccurrence(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProgramOccurrence(result)), nil
}
