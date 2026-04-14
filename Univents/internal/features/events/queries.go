package events

import (
	"context"
	"univents/internal/platform/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	events   ports.EventsRepository
	gaClient *goauth.Client
	tracer   trace.Tracer
	az       *authzed.Client
	tx       database.TxRunner
}

func NewQueryService(
	events ports.EventsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		events:   events,
		gaClient: gaClient,
		tracer:   tracer,
		az:       az,
		tx:       tx,
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
