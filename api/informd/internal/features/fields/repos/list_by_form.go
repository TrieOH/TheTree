package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListByFormID(ctx context.Context, formID uuid.UUID) ([]models.Field, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldRepo.ListByFormID")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListFieldsByFormID(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapField), nil
}
