package queries

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (q *Queries) ListOwned(ctx context.Context) ([]models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOwned")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	events, err := q.events.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
