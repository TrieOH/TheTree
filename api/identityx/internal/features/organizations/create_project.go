package organizations

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"
)

func (o *Operations) CreateProject(ctx context.Context, in models.CreateOrgProjectInput) (*models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "CreateProject")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := o.orgs.GetByID(ctx, in.OrganizationID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return nil, err
	}

	project, err := models.NewProject(org.OwnerID, in.BrandSlug, in.Name, in.Domain, &in.OrganizationID)
	if err != nil {
		return nil, err
	}

	return o.projects.Create(ctx, *project)
}
