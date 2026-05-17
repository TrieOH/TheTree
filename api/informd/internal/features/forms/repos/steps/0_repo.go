package steps

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type stepRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.StepRepo = (*stepRepo)(nil)

func NewStepRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.StepRepo {
	return &stepRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("step"),
	}
}

func mapStep(src sqlc.Step) models.Step {
	return models.Step{
		ID:          src.ID,
		FormID:      src.FormID,
		Title:       src.Title,
		Description: src.Description,
	}
}
