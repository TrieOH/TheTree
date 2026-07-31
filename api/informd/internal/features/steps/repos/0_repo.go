package repos

import (
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.StepRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("step"),
	}
}

func mapStep(src sqlc.Step) models.Step {
	return models.Step{
		ID:           src.ID,
		FormID:       src.FormID,
		Title:        src.Title,
		Description:  src.Description,
		PositionHint: src.PositionHint,
	}
}
