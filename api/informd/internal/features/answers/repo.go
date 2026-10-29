package answers

import (
	"Informd/errs"
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.AnswerRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrAnswerNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapAnswer(src sqlc.Answer) models.Answer {
	return models.Answer{
		ID:         src.ID,
		ResponseID: src.ResponseID,
		FieldID:    src.FieldID,
		Answer:     src.Answer,
		AnsweredAt: src.AnsweredAt,
	}
}

func (repo *Repo) BatchUpsert(ctx context.Context, answers []models.Answer) error {
	params := xslices.MapSlice(answers, models.ToBatchUpsertAnswersParams)
	return database.BatchExec(
		repo.db(ctx).BatchUpsertAnswers(ctx, params),
		repo.dbe,
		func(i int) string { return params[i].FieldID.String() },
	)
}

func (repo *Repo) GetByField(ctx context.Context, fieldID uuid.UUID) ([]models.Answer, error) {
	rows, err := repo.db(ctx).GetAnswersByField(ctx, &fieldID)
	return database.Many(rows, err, repo.dbe, mapAnswer)
}

func (repo *Repo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]models.Answer, error) {
	rows, err := repo.db(ctx).GetAnswersByFormID(ctx, formID)
	return database.Many(rows, err, repo.dbe, mapAnswer)
}

func (repo *Repo) GetByResponse(ctx context.Context, responseID uuid.UUID) ([]models.Answer, error) {
	rows, err := repo.db(ctx).GetAnswersByResponse(ctx, responseID)
	return database.Many(rows, err, repo.dbe, mapAnswer)
}
