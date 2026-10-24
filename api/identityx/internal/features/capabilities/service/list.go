package service

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	ctx, span := telemetry.StartSpan(ctx, "List")
	defer span.End()

	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Capabilities.List(ctx, projectID)
}
