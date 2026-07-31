package repos

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetMember(ctx context.Context, memberID, orgID uuid.UUID) (*models.OrganizationMember, error) {
	ctx, span := repo.tracer.Start(ctx, "GetMember")
	defer span.End()
	sqlcMember, err := database.Queries(ctx, repo.q).GetOrganizationMember(ctx, sqlc.GetOrganizationMemberParams{
		MemberID:       memberID,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOrganizationMember(sqlcMember)), nil
}
