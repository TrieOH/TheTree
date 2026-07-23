package queries

import (
	"context"
	idx "sdk/identityx"
	"univents/models"
)

func (q *Queries) ListJoined(ctx context.Context) ([]models.Event, error) {
	ctx, span := q.tracer.Start(ctx, "ListJoined")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	events, err := q.events.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return events, nil
}
