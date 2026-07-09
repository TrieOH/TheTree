package commands

import (
	"context"
	"lib/objectstorage"
	"univents/contracts"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) UnsetLogo(ctx context.Context, id uuid.UUID) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.UnsetLogo")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("unset_logo.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event.LogoUrl == nil {
		return nil, errx.Invalid("event").SetMessage("already has no logo")
	}

	bucket, key, err := objectstorage.ParseURL(*event.LogoUrl)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.obj.RemoveObject(ctx, bucket, key); err != nil {
		return nil, errx.Internal("event").SetMessage("failed to delete image from storage: " + err.Error())
	}

	event, err = uc.events.UnsetLogo(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	return event, nil
}
