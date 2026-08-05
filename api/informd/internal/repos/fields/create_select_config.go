package fields

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) CreateSelectConfig(ctx context.Context, toCreate models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldRepo.CreateSelectConfig")
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
