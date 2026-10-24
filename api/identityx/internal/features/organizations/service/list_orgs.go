package service

import (
	"IdentityX/models"
	"context"
)

func (o *Operations) ListOrgs(ctx context.Context) ([]models.Organization, error) {
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
