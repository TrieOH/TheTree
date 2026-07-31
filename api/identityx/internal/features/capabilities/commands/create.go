package commands

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateCapabilityInput) (*models.Capability, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = c.authz.CheckProject(ctx, ident.Sub.ID, *payload.ProjectID, nil, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	capability := models.Capability{
		ProjectID: payload.ProjectID,
		Resource:  payload.Resource,
		Action:    payload.Action,
		CreatedBy: ident.Sub.ID,
	}

	return c.capabilities.Create(ctx, capability)
}
