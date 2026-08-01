package events

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) ListJoined(ctx context.Context) ([]models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListJoined")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	events, err := o.events.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
