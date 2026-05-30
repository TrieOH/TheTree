package repos

import (
	"Informd/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetSelectConfig(ctx context.Context, fieldID uuid.UUID) (*models.FieldSelectConfig, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FieldRepo.GetSelectConfig")
	defer span.End()
	sqlcConfig, err := database.Queries(ctx, repo.q).GetFieldSelectConfig(ctx, fieldID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapFieldSelectConfig(sqlcConfig)), nil
}
