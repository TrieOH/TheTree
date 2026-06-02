package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.AnswerRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.AnswerRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("answer"),
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
