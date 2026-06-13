package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListByStepID(ctx context.Context, stepID uuid.UUID) ([]models.Field, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FieldRepo.ListByStepID")
	defer span.End()
	sqlcFields, err := database.Queries(ctx, repo.q).ListFieldsByStepID(ctx, stepID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcFields, mapField), nil
}
