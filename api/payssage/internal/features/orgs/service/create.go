package service

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) Create(ctx context.Context, in models.CreateOrganizationInput) (*models.Organization, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org := models.Organization{
		OwnerID: ident.Sub.ID,
		Name:    in.Name,
		Slug:    in.Slug,
	}

	var created *models.Organization
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Orgs.Create(ctx, org)
		if err != nil {
			return err
		}

		owner := models.OrganizationMember{
			MemberID:       ident.Sub.ID,
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
