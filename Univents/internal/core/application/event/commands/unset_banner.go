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

func (uc *CommandService) UnsetBanner(ctx context.Context, id uuid.UUID) (event *domain.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.UnsetBanner")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("unset_banner.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event.BannerUrl == nil {
		return nil, errx.Invalid("event").SetMessage("already has no banner")
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
		return nil, err
	}

	bucket, key, err := parseMinioURL(*event.BannerUrl)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.minio.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return nil, errx.Internal("event").SetMessage("failed to delete image from storage: " + err.Error())
	}

	event, err = uc.events.UnsetBanner(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	return event, nil
}
