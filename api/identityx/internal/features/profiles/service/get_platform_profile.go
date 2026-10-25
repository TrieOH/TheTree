package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

// GetPlatformProfile returns a platform-scoped actor's profile (project_id is NULL).
func (o *Operations) GetPlatformProfile(ctx context.Context, actorID uuid.UUID) (*models.ActorProfile, error) {
	profile, err := o.Profiles.Get(ctx, actorID)
	if err != nil {
		return nil, err
	}

	return o.migrateOnDemand(ctx, profile, nil)
}
