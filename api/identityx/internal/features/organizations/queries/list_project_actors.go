package queries

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (q *Queries) ListProjectActors(ctx context.Context, orgID, projectID uuid.UUID) ([]models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListProjectActors")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = q.authz.CheckProject(ctx, ident.Sub.ID, projectID, &orgID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return q.actors.List(ctx, projectID)
}
