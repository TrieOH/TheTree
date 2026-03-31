package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) SetLogo(ctx context.Context, id uuid.UUID, url string) (event *domain.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.SetLogo")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("set_logo.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.gaClient.Authz.Check().User(sub.ID).
		Object("events").
		Action("edit").
		Scope(event.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("event").SetMessage("insufficient permissions")
	}

	event, err = uc.events.SetLogo(ctx, event.ID, url)
	if err != nil {
		return nil, err
	}

	return event, nil
}
