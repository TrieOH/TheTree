package profiles

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) GetProfile(ctx context.Context, actorID, projectID uuid.UUID) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetProfile")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckProject(ctx, ident.Sub.ID, projectID, nil, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.profiles.Get(ctx, actorID)
}
