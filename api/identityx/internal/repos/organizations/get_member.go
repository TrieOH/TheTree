package organizations

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetMember(ctx context.Context, actorID, orgID uuid.UUID) (*models.OrganizationMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetMember")
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
