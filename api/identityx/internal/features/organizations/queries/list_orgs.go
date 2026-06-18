package queries

import (
	"IdentityX/models"
	"context"
)

func (q *Queries) ListOrgs(ctx context.Context) ([]models.Organization, error) {
	ctx, span := q.tracer.Start(ctx, "OrganizationService.ListOrgs")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownOrgs, err := q.orgs.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedOrgs, err := q.orgs.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownOrgs, joinedOrgs...), nil
}
