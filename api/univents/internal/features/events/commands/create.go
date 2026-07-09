package commands

import (
	"context"
	idx "sdk/identityx"
	"univents/contracts"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) CreateEvent(ctx context.Context, in contracts.CreateEventSpec) (out *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var validEvent *contracts.Event
	validEvent, err = contracts.NewEvent(ident.Sub.ID, &ident.Sub.ID, in)
	if err != nil {
		return nil, err
	}

	var created *contracts.Event
	created, err = uc.events.CreateEvent(ctx, validEvent)
	if err != nil {
		return nil, err
	}

	return created, nil
}
