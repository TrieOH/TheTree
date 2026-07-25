package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (q *Queries) ListPublic(ctx context.Context) ([]models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListPublic")
	defer span.End()

	events, err := q.events.ListPublic(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}
