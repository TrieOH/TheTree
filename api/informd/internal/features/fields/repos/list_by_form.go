package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListByFormID(ctx context.Context, formID uuid.UUID) ([]models.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "FieldRepo.ListByFormID")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListFieldsByFormID(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapField), nil
}
