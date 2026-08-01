package organizations

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"
)

func (o *Operations) Create(ctx context.Context, in models.CreateOrganizationInput) (*models.Organization, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := models.NewOrganization(ident.Sub.ID, in.Name, in.Slug)
	if err != nil {
		return nil, err
	}

	var created *models.Organization
	err = database.RunTx(ctx, func(ctx context.Context) error {
		created, err = o.orgs.Create(ctx, *org)
		if err != nil {
			return err
		}

		owner := models.OrganizationMember{
			ActorID:        ident.Sub.ID,
			OrganizationID: created.ID,
			Role:           models.OrganizationRoleOwner,
		}

		return o.orgs.AddMember(ctx, owner)
	})
	if err != nil {
		return nil, err
	}
	return created, err
}
