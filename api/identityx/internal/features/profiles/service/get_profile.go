package service

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) GetProfile(ctx context.Context, actorID, projectID uuid.UUID) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetProfile")
	defer span.End()

	// public read: the project route serves project-scoped actors only
	err := requireActorInProject(ctx, o.Actors, actorID, projectID)
	if err != nil {
		return nil, err
	}

	profile, err := o.Profiles.Get(ctx, actorID)
	if err != nil {
		return nil, err
	}

	return o.migrateOnDemand(ctx, profile, &projectID)
}
