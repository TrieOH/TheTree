package repo

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListByStepID(ctx context.Context, stepID uuid.UUID) ([]models.Field, error) {
	sqlcFields, err := database.Queries(ctx, repo.q).ListFieldsByStepID(ctx, stepID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcFields, mapField), nil
}
