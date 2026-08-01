package organizations

import (
	"IdentityX/internal/sqlc"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, actorID, orgID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
	defer span.End()
	err := database.Queries(ctx, repo.q).RemoveOrganizationMember(ctx, sqlc.RemoveOrganizationMemberParams{
		ActorID:        actorID,
		OrganizationID: orgID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
