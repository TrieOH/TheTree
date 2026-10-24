package service

import (
	"IdentityX/models"
	"context"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateCapabilityInput) (*models.Capability, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckProject(ctx, *payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	capability := models.Capability{
		ProjectID: payload.ProjectID,
		Resource:  payload.Resource,
		Action:    payload.Action,
		CreatedBy: ident.Sub.ID,
	}

	return o.Capabilities.Create(ctx, capability)
}
