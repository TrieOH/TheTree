package repo

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, memberID, orgID uuid.UUID) error {
	err := database.Queries(ctx, repo.q).RemoveOrganizationMember(ctx, sqlc.RemoveOrganizationMemberParams{
		MemberID:       memberID,
		OrganizationID: orgID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
