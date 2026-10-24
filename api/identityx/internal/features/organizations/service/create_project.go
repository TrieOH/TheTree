package service

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"
)

func (o *Operations) CreateProject(ctx context.Context, in models.CreateOrgProjectInput) (*models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "CreateProject")
	defer span.End()
	org, err := o.Orgs.GetByID(ctx, in.OrganizationID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckOrg(ctx, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return nil, err
	}

	project, err := models.NewProject(org.OwnerID, in.BrandSlug, in.Name, in.Domain, &in.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Project creation and key provisioning are one transaction: the
	// Key-lifecycle module provisions signing and encryption keys for the
	// new project here, so an org-created project is never token-broken
	// until the next boot.
	var created *models.Project
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Projects.Create(ctx, *project)
		if err != nil {
			return err
		}
		return o.Keys.Ensure(ctx, &created.ID)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
