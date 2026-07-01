package commands

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
)

func (c *Commands) CreateProject(ctx context.Context, in models.CreateOrgProjectInput) (*models.Project, error) {
	ctx, span := c.tracer.Start(ctx, "OrganizationService.CreateProject")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := c.orgs.GetByID(ctx, in.OrganizationID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != org.OwnerID {
		member, err := c.orgs.GetMember(ctx, ident.Sub.ID, org.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
		if member.Role != models.OrganizationRoleAdmin {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	project, err := models.NewProject(org.OwnerID, in.BrandSlug, in.Name, in.Domain, &in.OrganizationID)
	if err != nil {
		return nil, err
	}

	return c.projects.Create(ctx, *project)
}
