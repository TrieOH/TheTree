package commands

import (
	"context"
	"univents/internal/eventcore/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) Create(ctx context.Context, toCreate domain.Event) (out *domain.Event, err error) {
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

	allowed, err := ga.Authz.Check().User(sub.ID).
		Object(goauth.Object("events", "*")).
		Action(goauth.Action("create")).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fail.New(errx.AuthzInsufficientPermissions)
	}

	toCreate.CreatedBy = sub.ID

	var eventUUID uuid.UUID
	eventUUID, err = uuid.NewV7()
	if err != nil {
		return nil, fail.New(errx.SYSUUIDV7GenerationError).RecordCtx(ctx)
	}
	eventIDStr := eventUUID.String()
	toCreate.ID = eventUUID

	var scope *goauth.Scope
	scope, err = ga.Scopes.Create(ctx, toCreate.Slug, &eventIDStr)
	if err != nil {
		return nil, fail.AsFail(err).System().RecordCtx(ctx)
	}

	toCreate.GoauthScopeID = scope.ID

	var created *domain.Event
	created, err = uc.events.Create(ctx, toCreate) // FIXME if this fails the scope must be undone (SAGA PATTERN)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.String("event.id", eventIDStr))

	return created, nil
}
