package repo

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"
)

func (repo *Repo) BulkEdit(ctx context.Context, steps []models.Step) error {
	params := xslices.MapSlice(steps, models.ToBulkEditStepsParams)
	return database.BatchExec(
		database.Queries(ctx, repo.q).BulkEditSteps(ctx, params),
		repo.dbe,
		func(i int) string { return params[i].ID.String() },
	)
}
