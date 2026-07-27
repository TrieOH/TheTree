package queries

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
	"lib/telemetry"
)

// GetPlatformProfile returns a platform-scoped actor's profile (project_id is NULL).
func (q *Queries) GetPlatformProfile(ctx context.Context, actorID uuid.UUID) (*models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetPlatformProfile")
	defer span.End()

	_, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return q.profiles.Get(ctx, actorID)
}
