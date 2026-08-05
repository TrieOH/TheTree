package fields

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListByStepID(ctx context.Context, stepID uuid.UUID) ([]models.Field, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldRepo.ListByStepID")
	defer span.End()
	sqlcFields, err := database.Queries(ctx, repo.q).ListFieldsByStepID(ctx, stepID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcFields, mapField), nil
}
