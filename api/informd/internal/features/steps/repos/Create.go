package repos

import (
	"context"

	"Informd/internal/database/sqlc"
	"Informd/models"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, toCreate models.Step) (*models.Step, error) {
	ctx, span := repo.tracer.Start(ctx, "StepRepo.Create")
	defer span.End()
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
