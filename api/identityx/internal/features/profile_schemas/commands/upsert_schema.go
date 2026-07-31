package commands

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"
)

func (c *Commands) UpsertSchema(ctx context.Context, payload models.UpsertProfileSchemaInput) (*models.ProjectProfileSchema, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpsertSchema")
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

	err = c.authz.CheckProject(ctx, ident.Sub.ID, *payload.ProjectID, nil, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	return c.schemas.Upsert(ctx, models.ProjectProfileSchema{
		ProjectID: payload.ProjectID,
		Schema:    payload.Schema,
		Active:    payload.Active,
	})
}
