package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) GetProfile(ctx context.Context, actorID, projectID uuid.UUID) (*models.ActorProfile, error) {
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
