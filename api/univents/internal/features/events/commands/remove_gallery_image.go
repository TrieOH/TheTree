package commands

import (
	"context"
	"lib/objectstorage"
	"univents/contracts"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.RemoveGalleryImage")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("remove_gallery.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	bucket, key, err := objectstorage.ParseURL(url)
	if err != nil {
		return nil, errx.Invalid("event").SetMessage("invalid image url")
	}

	if err = uc.obj.RemoveObject(ctx, bucket, key); err != nil {
		return nil, errx.Internal("event").SetMessage("failed to delete image from storage: " + err.Error())
	}

	event, err = uc.events.RemoveGalleryImage(ctx, event.ID, url)
	if err != nil {
		return nil, err
	}

	return event, nil
}
