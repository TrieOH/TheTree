package queries

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (q *Queries) GetByEmail(ctx context.Context, email string, projectID uuid.UUID) (*models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByEmail")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = q.authz.CheckProject(ctx, ident.Sub.ID, projectID, nil, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	actor, err := q.actors.GetByEmail(ctx, email, &projectID)
	if err != nil {
		return nil, err
	}

	return actor, nil
}
