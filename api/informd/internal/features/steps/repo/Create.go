package repo

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Step) (*models.Step, error) {
	sqlcStep, err := database.Queries(ctx, repo.q).CreateStep(ctx, sqlc.CreateStepParams{
		FormID:       toCreate.FormID,
		Title:        toCreate.Title,
		Description:  toCreate.Description,
		PositionHint: toCreate.PositionHint,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapStep(sqlcStep)), nil
}
