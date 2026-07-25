package repos

import (
	sqlc2 "Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type repo struct {
	q      *sqlc2.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.AnswerRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.AnswerRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("answer"),
	}
}

func mapAnswer(src sqlc2.Answer) models.Answer {
	return models.Answer{
		ID:         src.ID,
		ResponseID: src.ResponseID,
		FieldID:    src.FieldID,
		Answer:     src.Answer,
		AnsweredAt: src.AnsweredAt,
	}
}
