package service

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) ListOrgs(ctx context.Context) ([]models.Organization, error) {
	ident, err := idx.RequireIdentity(ctx)
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
