package queries

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"
)

func (q *Queries) ListByProfile(ctx context.Context) ([]models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByProfile")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return q.intents.ListByOwner(ctx, ident.Sub.ID)
}
