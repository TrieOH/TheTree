package commands

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateCapabilityInput) (*models.Capability, error) {
	ctx, span := c.tracer.Start(ctx, "Create")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var project *models.Project
	project, err = c.projects.GetByID(ctx, *payload.ProjectID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != project.OwnerID {
		member, err := c.projects.GetMember(ctx, ident.Sub.ID, project.ID)
		if err != nil {
			return nil, err
		}
		if member.Role != models.ProjectRoleAdmin {
			return nil, fun.ErrForbidden("insufficient permissions")
		}
	}

	capability := models.Capability{
		ProjectID: &project.ID,
		Resource:  payload.Resource,
		Action:    payload.Action,
		CreatedBy: ident.Sub.ID,
	}

	return c.capabilities.Create(ctx, capability)
}
