package queries

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (q *Queries) GetProfile(ctx context.Context, actorID, projectID uuid.UUID) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetProfile")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	project, err := q.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != project.OwnerID {
		_, err = q.projects.GetMember(ctx, ident.Sub.ID, projectID)
		if err != nil {
			return nil, err
		}
	}

	return q.profiles.Get(ctx, actorID)
}
