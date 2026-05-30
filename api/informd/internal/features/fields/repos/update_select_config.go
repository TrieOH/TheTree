package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"
)

func (repo *repo) UpdateSelectConfig(ctx context.Context, toUpdate models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FieldRepo.UpdateSelectConfig")
	defer span.End()
	sqlcConfig, err := database.Queries(ctx, repo.q).UpdateFieldSelectConfig(ctx, sqlc.UpdateFieldSelectConfigParams{
		FieldID:   toUpdate.FieldID,
		Behaviour: string(toUpdate.Behaviour),
		ValueType: string(toUpdate.ValueType),
		Options:   toUpdate.Options,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapFieldSelectConfig(sqlcConfig)), nil
}
