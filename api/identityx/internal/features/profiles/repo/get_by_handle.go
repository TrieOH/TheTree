package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
)

// GetByHandle returns the profile whose handle matches. Handles are unique
// when present, so at most one profile can match.
func (r *Repo) GetByHandle(ctx context.Context, handle string) (*models.ActorProfile, error) {
	result, err := database.Queries(ctx, r.q).GetActorProfileByHandle(ctx, &handle)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapActorProfile(result)), nil
}
