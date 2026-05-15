package steps

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/errx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type stepRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    *errx.DBHandler
}

var _ ports.StepRepo = (*stepRepo)(nil)

func NewStepRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer, dbe *errx.DBHandler) ports.StepRepo {
	return &stepRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    dbe,
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
