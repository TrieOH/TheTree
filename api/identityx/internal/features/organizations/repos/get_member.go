package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetMember(ctx context.Context, actorID, orgID uuid.UUID) (*models.OrganizationMember, error) {
	ctx, span := repo.tracer.Start(ctx, "GetMember")
	defer span.End()
	sqlcMember, err := database.Queries(ctx, repo.q).GetOrganizationMember(ctx, sqlc.GetOrganizationMemberParams{
		ActorID:        actorID,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOrganizationMember(sqlcMember)), nil
}
