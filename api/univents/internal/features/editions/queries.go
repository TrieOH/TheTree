package editions

import (
	"context"

	"lib/database"
	"univents/contracts"
	"univents/internal/shared/ports"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueryService struct {
	events   ports.EventsRepository
	editions ports.EditionsRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewQueryService(
	events ports.EventsRepository,
	editions ports.EditionsRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		events:   events,
		editions: editions,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	}
}

func (uc *QueryService) ListEditions(ctx context.Context, eventID uuid.UUID) (out []contracts.Edition, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EditionsService.ListEditions")
	defer span.End()

	var outEditions []contracts.Edition
	outEditions, err = uc.editions.List(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return outEditions, nil
}

func (uc *QueryService) ListEditionsAdmin(ctx context.Context, eventID uuid.UUID) (out []contracts.Edition, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EditionsService.ListEditionsAdmin")
	defer span.End()

	var outEditions []contracts.Edition
	outEditions, err = uc.editions.ListAdmin(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return outEditions, nil
}
