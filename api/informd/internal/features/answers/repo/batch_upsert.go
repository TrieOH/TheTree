package repo

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"
)

func (repo *Repo) BatchUpsert(ctx context.Context, answers []models.Answer) error {
	params := xslices.MapSlice(answers, models.ToBatchUpsertAnswersParams)
	return database.NamedBatchExec(
		database.Queries(ctx, repo.q).BatchUpsertAnswers(ctx, params),
		repo.dbe,
		func(i int) string { return params[i].FieldID.String() },
	)
}
