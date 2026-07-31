package queries

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"
)

func (q *Queries) ListOwned(ctx context.Context) ([]models.Collector, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOwned")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return q.collectors.ListByOwner(ctx, ident.Sub.ID)
}
