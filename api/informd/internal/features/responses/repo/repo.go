package repo

import (
	"Informd/errs"
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ResponseRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrResponseNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapResponse(src sqlc.Response) models.Response {
	return models.Response{
		ID:          src.ID,
		FormID:      src.FormID,
		InviteID:    src.InviteID,
		ResponderID: src.ResponderID,
		Email:       src.Email,
		StartedAt:   src.StartedAt,
		FinishedAt:  src.FinishedAt,
	}
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Response) (*models.Response, error) {
	row, err := repo.db(ctx).CreateResponse(ctx, sqlc.CreateResponseParams{
		FormID:      toCreate.FormID,
		InviteID:    toCreate.InviteID,
		ResponderID: toCreate.ResponderID,
		Email:       toCreate.Email,
	})
	return database.One(row, err, repo.dbe, mapResponse)
}

func (repo *Repo) Finish(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).FinishResponse(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Response, error) {
	row, err := repo.db(ctx).GetResponseByID(ctx, id)
	return database.One(row, err, repo.dbe, mapResponse)
}

func (repo *Repo) ListByForm(ctx context.Context, formID uuid.UUID) ([]models.Response, error) {
	rows, err := repo.db(ctx).ListResponsesByForm(ctx, formID)
	return database.Many(rows, err, repo.dbe, mapResponse)
}
