package events

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) ListOwned(ctx context.Context) ([]models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOwned")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	events, err := o.events.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
