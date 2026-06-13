package repos

import (
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Field, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FieldRepo.GetByID")
	defer span.End()
	sqlcField, err := database.Queries(ctx, repo.q).GetFieldByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapField(sqlcField)), nil
}
