package service

import (
	"IdentityX/models"
	"context"
)

// GetProfileByHandle returns the profile with the given handle. Handles are
// globally unique when present, so at most one profile can match. Public
// read, like the other profile get routes.
func (o *Operations) GetProfileByHandle(ctx context.Context, handle string) (*models.ActorProfile, error) {
	profile, err := o.Profiles.GetByHandle(ctx, handle)
	if err != nil {
		return nil, err
	}

	// migrate against the actor's scope (project or platform)
	actor, err := o.Actors.GetByID(ctx, profile.ActorID)
	if err != nil {
		return nil, err
	}

	return o.migrateOnDemand(ctx, profile, actor.ProjectID)
}
