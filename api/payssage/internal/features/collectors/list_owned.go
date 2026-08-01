package collectors

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) ListOwned(ctx context.Context) ([]models.Collector, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOwned")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return o.collectors.ListByOwner(ctx, ident.Sub.ID)
}
