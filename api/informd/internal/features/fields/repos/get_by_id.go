package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Field, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldRepo.GetByID")
	defer span.End()
	sqlcField, err := database.Queries(ctx, repo.q).GetFieldByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapField(sqlcField)), nil
}
