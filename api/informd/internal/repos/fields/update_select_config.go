package fields

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) UpdateSelectConfig(ctx context.Context, toUpdate models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldRepo.UpdateSelectConfig")
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
