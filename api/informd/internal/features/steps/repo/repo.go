package repo

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

var _ ports.StepRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrStepNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapStep(src sqlc.Step) models.Step {
	return models.Step{
		ID:           src.ID,
		FormID:       src.FormID,
		Title:        src.Title,
		Description:  src.Description,
		PositionHint: src.PositionHint,
	}
}

func (repo *Repo) BulkEdit(ctx context.Context, steps []models.Step) error {
	params := xslices.MapSlice(steps, models.ToBulkEditStepsParams)
	return database.BatchExec(
		repo.db(ctx).BulkEditSteps(ctx, params),
		repo.dbe,
		func(i int) string { return params[i].ID.String() },
	)
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Step) (*models.Step, error) {
	sqlcStep, err := repo.db(ctx).CreateStep(ctx, sqlc.CreateStepParams{
		FormID:       toCreate.FormID,
		Title:        toCreate.Title,
		Description:  toCreate.Description,
		PositionHint: toCreate.PositionHint,
	})
	return database.One(sqlcStep, err, repo.dbe, mapStep)
}

func (repo *Repo) List(ctx context.Context, formID uuid.UUID) ([]models.Step, error) {
	sqlcForm, err := repo.db(ctx).ListStepsByFormID(ctx, formID)
	return database.Many(sqlcForm, err, repo.dbe, mapStep)
}
