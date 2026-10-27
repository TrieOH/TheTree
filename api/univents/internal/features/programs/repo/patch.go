package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) Patch(ctx context.Context, id uuid.UUID, program *models.Program) (*models.Program, error) {
	result, err := database.Queries(ctx, repo.q).PatchProgram(ctx, sqlc.PatchProgramParams{
		Kind:           string(program.Kind),
		Name:           program.Name,
		Description:    program.Description,
		MinAccessLevel: program.MinAccessLevel,
		StaffOnly:      program.StaffOnly,
		Price:          priceValue(program.Price),
		BannerUrl:      program.BannerURL,
		ID:             id,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProgram(result)), nil
}

func (repo *Repo) PatchOccurrence(ctx context.Context, id uuid.UUID, occurrence *models.ProgramOccurrence) (*models.ProgramOccurrence, error) {
	result, err := database.Queries(ctx, repo.q).PatchProgramOccurrence(ctx, sqlc.PatchProgramOccurrenceParams{
		StartsAt:    occurrence.StartsAt,
		EndsAt:      occurrence.EndsAt,
		MaxCapacity: occurrence.MaxCapacity,
		ID:          id,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProgramOccurrence(result)), nil
}
