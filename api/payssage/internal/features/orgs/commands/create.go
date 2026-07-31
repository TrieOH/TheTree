package commands

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"
)

func (c *Commands) Create(ctx context.Context, in models.CreateOrganizationInput) (*models.Organization, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

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
	err = database.RunTx(ctx, func(ctx context.Context) error {
		created, err = c.orgs.Create(ctx, org)
		if err != nil {
			return err
		}

		owner := models.OrganizationMember{
			MemberID:       ident.Sub.ID,
			OrganizationID: created.ID,
			Role:           models.OrganizationRoleOwner,
		}

		return c.orgs.AddMember(ctx, owner)
	})
	if err != nil {
		return nil, err
	}

	return created, err
}
