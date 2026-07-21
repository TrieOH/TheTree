package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Collector, error) {
	ctx, span := q.tracer.Start(ctx, "ListByOrg")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := q.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if org.OwnerID == ident.Sub.ID {
		return q.collectors.ListByOrg(ctx, org.ID)
	}

	_, err = q.orgs.GetMember(ctx, ident.Sub.ID, org.ID)
	if err != nil {
		return nil, err
	}

	return q.collectors.ListByOrg(ctx, org.ID)
}
