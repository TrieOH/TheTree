package repo

import (
	"Informd/errs"
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.NamedErrorHandler
}

var _ ports.ResponseRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewNamedErrorHandler(errs.ErrResponseNotFound),
	}
}

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
