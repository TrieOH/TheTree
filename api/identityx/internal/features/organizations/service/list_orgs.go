package service

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"
)

func (o *Operations) ListOrgs(ctx context.Context) ([]models.Organization, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOrgs")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	ownOrgs, err := o.Orgs.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	joinedOrgs, err := o.Orgs.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return append(ownOrgs, joinedOrgs...), nil
}
