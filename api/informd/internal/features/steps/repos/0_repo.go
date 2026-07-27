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

var _ ports.StepRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.StepRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("step"),
	}
}

func mapStep(src sqlc2.Step) models.Step {
	return models.Step{
		ID:           src.ID,
		FormID:       src.FormID,
		Title:        src.Title,
		Description:  src.Description,
		PositionHint: src.PositionHint,
	}
}
