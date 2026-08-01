package orgs

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, memberID, orgID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
	defer span.End()
	err := database.Queries(ctx, repo.q).RemoveOrganizationMember(ctx, sqlc.RemoveOrganizationMemberParams{
		MemberID:       memberID,
		OrganizationID: orgID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
