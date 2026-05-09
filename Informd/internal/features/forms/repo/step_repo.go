package repo

import (
	"Informd/internal/platform/database"
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/errx"
	"Informd/internal/shared/ports"
	"Informd/internal/shared/xslices"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type stepRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.StepRepo = (*stepRepo)(nil)

func NewStepRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.StepRepo {
	return &stepRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func mapStep(src sqlc.Step) contracts.Step {
	return contracts.Step{
		ID:          src.ID,
		FormID:      src.FormID,
		Title:       src.Title,
		Description: src.Description,
	}
}

func (repo *stepRepo) Create(ctx context.Context, toCreate contracts.Step) (*contracts.Step, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Create")
	defer span.End()
	sqlcStep, err := database.Queries(ctx, repo.q).CreateStep(ctx, sqlc.CreateStepParams{
		FormID:       toCreate.FormID,
		Title:        toCreate.Title,
		Description:  toCreate.Description,
		PositionHint: toCreate.PositionHint,
	})
	if err != nil {
		return nil, errx.DB(err, "form")
	}
	return new(mapStep(sqlcStep)), nil
}

func (repo *stepRepo) List(ctx context.Context, formID uuid.UUID) ([]contracts.Step, error) {
	ctx, span := database.Span(ctx, repo.tracer, "List")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).ListStepsByFormID(ctx, formID)
	if err != nil {
		return nil, errx.DB(err, "form")
	}
	return xslices.MapSlice(sqlcForm, mapStep), nil
}
