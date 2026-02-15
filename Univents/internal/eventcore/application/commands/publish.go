package commands

import (
	"context"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
)

func (uc *CommandService) Publish(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "EventService.Publish")
	defer span.End()

	ga := uc.gaClient

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	event, err := uc.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	allowed, err := ga.Authz.Check().User(sub.ID).
		Object(goauth.Object("events", "*")).
		Action(goauth.Action("publish")).
		Scope(event.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return fail.New(errx.AuthzInsufficientPermissions)
	}

	err = uc.events.Publish(ctx, eventID)
	if err != nil {
		return err
	}

	return nil
}
