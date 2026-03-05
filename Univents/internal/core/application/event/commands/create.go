package commands

import (
	"context"
	"encoding/json"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) CreateEvent(ctx context.Context, in domain.CreateEventSpec) (out *domain.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	ga := uc.gaClient

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

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("events").
		Action("create").
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("event").SetMessage("insufficient permissions")
	}

	span.SetAttributes(attribute.String("event.id", validEvent.ID.String()))

	meta := json.RawMessage(`{"color": "#ae20fa", "icon": "CalendarRange"}`)
	var scope *goauth.Scope
	var idStr = validEvent.ID.String()
	scope, err = ga.Scopes.Create(ctx, validEvent.Slug, &idStr, meta)
	if err != nil {
		return nil, err
	}
	validEvent.AddScope(scope.ID)

	err = uc.gaClient.Roles.Give(ctx, sub.ID, "Event Owner", &scope.ID)
	if err != nil {
		return nil, err
	}

	var created *domain.Event
	created, err = uc.events.CreateEvent(ctx, validEvent) // FIXME if this fails the scope must be undone (SAGA PATTERN)
	if err != nil {
		return nil, err
	}

	return created, nil
}
