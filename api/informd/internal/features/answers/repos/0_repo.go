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

var _ ports.AnswerRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("answer"),
	}
}

func mapAnswer(src sqlc.Answer) models.Answer {
	return models.Answer{
		ID:         src.ID,
		ResponseID: src.ResponseID,
		FieldID:    src.FieldID,
		Answer:     src.Answer,
		AnsweredAt: src.AnsweredAt,
	}
}
