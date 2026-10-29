package responders

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

var _ ports.ResponderRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrResponderNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapResponder(src sqlc.Responder) models.Responder {
	return models.Responder{
		ID:     src.ID,
		UserID: src.UserID,
		Email:  src.Email,
	}
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Responder) (*models.Responder, error) {
	row, err := repo.db(ctx).CreateResponder(ctx, sqlc.CreateResponderParams{
		UserID: toCreate.UserID,
		Email:  toCreate.Email,
	})
	return database.One(row, err, repo.dbe, mapResponder)
}

func (repo *Repo) GetByEmail(ctx context.Context, email string) (*models.Responder, error) {
	row, err := repo.db(ctx).GetResponderByEmail(ctx, email)
	return database.One(row, err, repo.dbe, mapResponder)
}

func (repo *Repo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]models.Responder, error) {
	rows, err := repo.db(ctx).GetRespondersByFormID(ctx, formID)
	return database.Many(rows, err, repo.dbe, mapResponder)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Responder, error) {
	row, err := repo.db(ctx).GetResponderByID(ctx, id)
	return database.One(row, err, repo.dbe, mapResponder)
}
