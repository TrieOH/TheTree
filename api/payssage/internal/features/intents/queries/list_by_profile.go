package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"
)

func (q *Queries) ListByProfile(ctx context.Context) ([]models.Intent, error) {
	ctx, span := q.tracer.Start(ctx, "ListByProfile")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return q.intents.ListByOwner(ctx, ident.Sub.ID)
}
