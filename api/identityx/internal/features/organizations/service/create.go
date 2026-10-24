package service

import (
	"IdentityX/models"
	"context"
)

func (o *Operations) Create(ctx context.Context, in models.CreateOrganizationInput) (*models.Organization, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := models.NewOrganization(ident.Sub.ID, in.Name, in.Slug)
	if err != nil {
		return nil, err
	}

	var created *models.Organization
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Orgs.Create(ctx, *org)
		if err != nil {
			return err
		}

		owner := models.OrganizationMember{
			ActorID:        ident.Sub.ID,
			OrganizationID: created.ID,
			Role:           models.OrganizationRoleOwner,
		}

		return o.Orgs.AddMember(ctx, owner)
	})
	if err != nil {
		return nil, err
	}
	return created, err
}
