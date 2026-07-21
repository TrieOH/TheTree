package commands

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
)

func (c *Commands) UpsertSchema(ctx context.Context, payload models.UpsertProfileSchemaInput) (*models.ProjectProfileSchema, error) {
	ctx, span := c.tracer.Start(ctx, "UpsertProfileSchema")
	defer span.End()

	// fixme: platform schema: any authenticated actor can set it for now
	// (in the future, restrict to platform admins)
	if payload.ProjectID == nil {
		_, err := models.RequireIdentity(ctx)
		if err != nil {
			return nil, err
		}
		return c.schemas.Upsert(ctx, models.ProjectProfileSchema{
			ProjectID: nil,
			Schema:    payload.Schema,
			Active:    payload.Active,
		})
	}

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	project, err := c.projects.GetByID(ctx, *payload.ProjectID)
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

	return c.schemas.Upsert(ctx, models.ProjectProfileSchema{
		ProjectID: payload.ProjectID,
		Schema:    payload.Schema,
		Active:    payload.Active,
	})
}
