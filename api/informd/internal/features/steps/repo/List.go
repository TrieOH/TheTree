package repo

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) List(ctx context.Context, formID uuid.UUID) ([]models.Step, error) {
	sqlcForm, err := database.Queries(ctx, repo.q).ListStepsByFormID(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForm, mapStep), nil
}
