package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"
)

func (repo *repo) Create(ctx context.Context, toCreate *models.Edition) (*models.Edition, error) {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.Create")
	defer span.End()
	edition, err := database.Queries(ctx, repo.q).CreateEdition(ctx, sqlc.CreateEditionParams{
		EventID:     toCreate.EventID,
		EditionName: toCreate.Name,
		Slug:        toCreate.Slug,
		StartsAt:    toCreate.StartsAt,
		EndsAt:      toCreate.EndsAt,
		CreatedBy:   toCreate.CreatedBy,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEdition(edition)), nil
}
