package commands

import (
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) UnsetLogo(ctx context.Context, id uuid.UUID) (event *domain.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.UnsetLogo")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("unset_logo.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event.LogoUrl == nil {
		return nil, errx.Invalid("event").SetMessage("already has no logo")
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

	bucket, key, err := parseMinioURL(*event.LogoUrl)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.minio.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return nil, errx.Internal("event").SetMessage("failed to delete image from storage: " + err.Error())
	}

	event, err = uc.events.UnsetLogo(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	return event, nil
}
