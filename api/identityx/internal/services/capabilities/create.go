package capabilities

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateCapabilityInput) (*models.Capability, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckProject(ctx, *payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	capability := models.Capability{
		ProjectID: payload.ProjectID,
		Resource:  payload.Resource,
		Action:    payload.Action,
		CreatedBy: ident.Sub.ID,
	}

	return o.capabilities.Create(ctx, capability)
}
