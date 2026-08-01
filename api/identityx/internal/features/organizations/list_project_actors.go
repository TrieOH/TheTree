package organizations

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) ListProjectActors(ctx context.Context, orgID, projectID uuid.UUID) ([]models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListProjectActors")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckProject(ctx, ident.Sub.ID, projectID, &orgID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.actors.List(ctx, projectID)
}
