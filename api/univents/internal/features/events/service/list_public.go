package service

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) ListPublic(ctx context.Context) ([]models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListPublic")
	defer span.End()

	events, err := o.Events.ListPublic(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}
