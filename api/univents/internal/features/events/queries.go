package events

import (
	"context"

	"lib/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type QueryService struct {
	events ports.EventsRepository
	logger *zap.Logger
	tracer trace.Tracer
	tx     database.TxRunner
}

func NewQueryService(
	events ports.EventsRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		events: events,
		logger: logger,
		tracer: tracer,
		tx:     tx,
	}
}

func (uc *QueryService) ListEvents(ctx context.Context) (out []contracts.Event, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EventService.ListEvents")
	defer span.End()

	var outEvents []contracts.Event
	outEvents, err = uc.events.ListEvents(ctx)
	if err != nil {
		return nil, err
	}

	return outEvents, nil
}

func (uc *QueryService) ListOwnEvents(ctx context.Context) (out []contracts.Event, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "EventService.ListOwnEvents")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var outEvents []contracts.Event
	outEvents, err = uc.events.ListOwnEvents(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return outEvents, nil
}
