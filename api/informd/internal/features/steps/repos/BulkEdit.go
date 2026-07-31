package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
)

func (repo *Repo) BulkEdit(ctx context.Context, steps []models.Step) error {
	ctx, span := telemetry.StartSpan(ctx, "StepRepo.BulkEdit")
	defer span.End()
	params := xslices.MapSlice(steps, models.ToBulkEditStepsParams)
	return database.BatchExec(
		database.Queries(ctx, repo.q).BulkEditSteps(ctx, params),
		repo.dbe,
		func(i int) string { return params[i].ID.String() },
	)
}
