package capabilities

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	ctx, span := telemetry.StartSpan(ctx, "List")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckProject(ctx, ident.Sub.ID, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.capabilities.List(ctx, projectID)
}
