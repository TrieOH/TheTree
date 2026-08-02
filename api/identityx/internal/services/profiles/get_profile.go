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

	// project users can read their own profile; any other read requires a
	// project member (member role or above)
	if ident.Sub.ID != actorID {
		err = o.authz.CheckProject(ctx, ident.Sub.ID, projectID, nil, models.ProjectRoleMember)
		if err != nil {
			return nil, err
		}
	}

	err = requireActorInProject(ctx, o.actors, actorID, projectID)
	if err != nil {
		return nil, err
	}

	profile, err := o.profiles.Get(ctx, actorID)
	if err != nil {
		return nil, err
	}

	return o.migrateOnDemand(ctx, profile, &projectID)
}
