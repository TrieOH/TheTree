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
	dbe database.ErrorHandler
}

var _ ports.ResponderRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrResponderNotFound),
	}
}

func mapResponder(src sqlc.Responder) models.Responder {
	return models.Responder{
		ID:     src.ID,
		UserID: src.UserID,
		Email:  src.Email,
	}
}
