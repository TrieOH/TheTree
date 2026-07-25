package queries

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
	"lib/telemetry"
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

	project, err := q.projects.GetByID(ctx, *projectID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != project.OwnerID {
		_, err = q.projects.GetMember(ctx, ident.Sub.ID, *projectID)
		if err != nil {
			return nil, err
		}
	}

	return q.schemas.Get(ctx, projectID)
}
