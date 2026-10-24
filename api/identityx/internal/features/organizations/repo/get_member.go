package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetMember(ctx context.Context, actorID, orgID uuid.UUID) (*models.OrganizationMember, error) {
	sqlcMember, err := database.Queries(ctx, repo.q).GetOrganizationMember(ctx, sqlc.GetOrganizationMemberParams{
		ActorID:        actorID,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOrganizationMember(sqlcMember)), nil
}
