package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) CreateEvent(ctx context.Context, in domain.CreateEventSpec) (out *domain.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var validEvent *domain.Event
	validEvent, err = domain.NewEvent(sub.ID, &sub.ID, in)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_events"),
		authz.Resource("platform", "global"),
	); err != nil {
		return nil, err
	}

	if err = authz.GrantRole(ctx, uc.az, "platform:global#event_creator@user:"+sub.ID.String()); err != nil {
		return nil, err
	} // FIXME Outbox this too

	var created *domain.Event
	created, err = uc.events.CreateEvent(ctx, validEvent) // FIXME if this fails the scope must be undone (SAGA PATTERN)
	if err != nil {
		return nil, err
	}

	return created, nil
}
