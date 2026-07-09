package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"
)

func (repo *activitiesRepo) Create(ctx context.Context, toCreate *contracts.Activity) (*contracts.Activity, error) {
	ctx, span := repo.tracer.Start(ctx, "ActivitiesRepo.Create")
	defer span.End()
	sqlcActivity, err := database.Queries(ctx, repo.q).CreateActivity(ctx, sqlc.CreateActivityParams{
		ID:                toCreate.ID,
		EditionID:         toCreate.EditionID,
		Title:             toCreate.Title,
		Description:       toCreate.Description,
		Location:          toCreate.Location,
		StartsAt:          toCreate.StartsAt,
		EndsAt:            toCreate.EndsAt,
		PresenterName:     toCreate.PresenterName,
		TokenCost:         toCreate.TokenCost,
		HasCapacity:       toCreate.HasCapacity,
		Capacity:          toCreate.Capacity,
		RemainingCapacity: toCreate.RemainingCapacity,
		Difficulty:        toCreate.Difficulty,
		CreatedBy:         toCreate.CreatedBy,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return mapActivityFromDB(&sqlcActivity), nil
}
