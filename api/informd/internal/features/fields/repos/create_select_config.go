package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"
)

func (repo *repo) CreateSelectConfig(ctx context.Context, toCreate models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FieldRepo.CreateSelectConfig")
	defer span.End()
	sqlcConfig, err := database.Queries(ctx, repo.q).CreateFieldSelectConfig(ctx, sqlc.CreateFieldSelectConfigParams{
		FieldID:   toCreate.FieldID,
		Behaviour: string(toCreate.Behaviour),
		ValueType: string(toCreate.ValueType),
		Options:   toCreate.Options,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapFieldSelectConfig(sqlcConfig)), nil
}
