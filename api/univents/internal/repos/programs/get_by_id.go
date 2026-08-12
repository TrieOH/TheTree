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

// GetByIDForUpdate is the row-lock variant used inside the checkout tx
// (split 7): serializes concurrent checkouts on the same program (price
// read) before availability is checked.
func (repo *Repo) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.GetByIDForUpdate")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetProgramByIDForUpdate(ctx, id)
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

// GetOccurrenceByIDForUpdate is the row-lock variant used inside the
// checkout tx (split 7): serializes concurrent checkouts on the same
// occurrence (capacity) before availability is checked.
func (repo *Repo) GetOccurrenceByIDForUpdate(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsRepo.GetOccurrenceByIDForUpdate")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetProgramOccurrenceByIDForUpdate(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProgramOccurrence(result)), nil
}
