package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

// GetRole returns the actor's platform role. A NotFound surfaces for actors
// with no platform role — callers map it (authz.CheckPlatform forbids).
func (repo *Repo) GetRole(ctx context.Context, actorID uuid.UUID) (models.PlatformRole, error) {
	role, err := database.Queries(ctx, repo.q).GetPlatformRoleByActor(ctx, actorID)
	if err != nil {
		return "", repo.dbe(err)
	}
	return models.PlatformRole(role), nil
}
