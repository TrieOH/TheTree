package editions

import (
	"context"
	"lib/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	events   ports.EventsRepository
	editions ports.EditionsRepository
	tracer   trace.Tracer
	az       *authzed.Client
	tx       database.TxRunner
}

func NewQueryService(
	events ports.EventsRepository,
	editions ports.EditionsRepository,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		events:   events,
		editions: editions,
		tracer:   tracer,
		az:       az,
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

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("view_editions"),
		authz.Resource("event", eventID.String()),
	); err != nil {
		return nil, err
	}

	var outEditions []contracts.Edition
	outEditions, err = uc.editions.ListAdmin(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return outEditions, nil
}
