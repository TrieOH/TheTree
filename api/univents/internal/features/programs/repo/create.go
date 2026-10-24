package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate *models.Program) (*models.Program, error) {
	result, err := database.Queries(ctx, repo.q).CreateProgram(ctx, sqlc.CreateProgramParams{
		EditionID:      toCreate.EditionID,
		Kind:           string(toCreate.Kind),
		Name:           toCreate.Name,
		Description:    toCreate.Description,
		MinAccessLevel: toCreate.MinAccessLevel,
		StaffOnly:      toCreate.StaffOnly,
		Price:          priceValue(toCreate.Price),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProgram(result)), nil
}

func (repo *Repo) CreateOccurrence(ctx context.Context, toCreate *models.ProgramOccurrence) (*models.ProgramOccurrence, error) {
	result, err := database.Queries(ctx, repo.q).CreateProgramOccurrence(ctx, sqlc.CreateProgramOccurrenceParams{
		ProgramID:   toCreate.ProgramID,
		EditionID:   toCreate.EditionID,
		StartsAt:    toCreate.StartsAt,
		EndsAt:      toCreate.EndsAt,
		MaxCapacity: toCreate.MaxCapacity,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProgramOccurrence(result)), nil
}
