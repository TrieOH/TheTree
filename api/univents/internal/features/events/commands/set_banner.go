package commands

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func (uc *Commands) SetBanner(ctx context.Context, id uuid.UUID, url string) (event *contracts.Event, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.SetBanner")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("set_banner.success", err == nil)) }()

	event, err = uc.events.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	event, err = uc.events.SetBanner(ctx, event.ID, url)
	if err != nil {
		return nil, err
	}

	return event, nil
}
