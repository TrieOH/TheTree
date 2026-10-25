package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Capabilities.List(ctx, projectID)
}
