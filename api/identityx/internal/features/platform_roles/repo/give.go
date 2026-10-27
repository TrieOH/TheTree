package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"encoding/json"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) Give(ctx context.Context, actorID uuid.UUID, role models.PlatformRole, metadata *json.RawMessage) (*models.PlatformRoleRelation, error) {
	sqlcPlatformRole, err := database.Queries(ctx, repo.q).GivePlatformRole(ctx, sqlc.GivePlatformRoleParams{
		ActorID:  actorID,
		Role:     string(role),
		Metadata: metadata,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapPlatformRole(sqlcPlatformRole)), nil
}
