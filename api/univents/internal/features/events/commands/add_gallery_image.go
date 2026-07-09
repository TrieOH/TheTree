package commands

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) AddGalleryImage(ctx context.Context, id uuid.UUID, url string) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.AddGalleryImage")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("add_gallery.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.AddGalleryImage(ctx, event.ID, url)
	if err != nil {
		return nil, err
	}

	return event, nil
}
