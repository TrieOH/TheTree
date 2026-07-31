package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
)

func (repo *Repo) BatchUpsert(ctx context.Context, answers []models.Answer) error {
	ctx, span := telemetry.StartSpan(ctx, "AnswerRepo.BatchUpsert")
	defer span.End()
	params := xslices.MapSlice(answers, models.ToBatchUpsertAnswersParams)
	return database.BatchExec(
		database.Queries(ctx, repo.q).BatchUpsertAnswers(ctx, params),
		repo.dbe,
		func(i int) string { return params[i].FieldID.String() },
	)
}
