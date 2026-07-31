package queries

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (q *Queries) List(ctx context.Context, projectID uuid.UUID) ([]models.Capability, error) {
	ctx, span := telemetry.StartSpan(ctx, "List")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = q.authz.CheckProject(ctx, ident.Sub.ID, projectID, nil, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return q.capabilities.List(ctx, projectID)
}
