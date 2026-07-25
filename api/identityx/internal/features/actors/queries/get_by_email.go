package queries

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
	"lib/telemetry"
)

func (q *Queries) GetByEmail(ctx context.Context, email string, projectID uuid.UUID) (*models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByEmail")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var project *models.Project
	project, err = q.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != project.OwnerID {
		_, err = q.projects.GetMember(ctx, ident.Sub.ID, projectID)
		if err != nil {
			return nil, err
		}
	}

	actor, err := q.actors.GetByEmail(ctx, email, &project.ID)
	if err != nil {
		return nil, err
	}

	return actor, nil
}
