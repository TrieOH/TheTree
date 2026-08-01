package programs

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.ListByEdition")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListProgramsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapProgram), nil
}

func (repo *Repo) ListOccurrencesByProgram(ctx context.Context, programID uuid.UUID) ([]models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.ListOccurrencesByProgram")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListProgramOccurrencesByProgram(ctx, programID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapProgramOccurrence), nil
}

func (repo *Repo) ListOccurrencesByEdition(ctx context.Context, editionID uuid.UUID) ([]models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.ListOccurrencesByEdition")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListProgramOccurrencesByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapProgramOccurrence), nil
}
