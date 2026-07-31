package queries

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (q *Queries) GetSchema(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetSchema")
	defer span.End()

	// platform schema: any authenticated actor can read it
	if projectID == nil {
		_, err := models.RequireIdentity(ctx)
		if err != nil {
			return nil, err
		}
		return q.schemas.Get(ctx, nil)
	}

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = q.authz.CheckProject(ctx, ident.Sub.ID, *projectID, nil, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return q.schemas.Get(ctx, projectID)
}
