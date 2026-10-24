package repo

import (
	"IdentityX/internal/sqlc"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, actorID, orgID uuid.UUID) error {
	err := database.Queries(ctx, repo.q).RemoveOrganizationMember(ctx, sqlc.RemoveOrganizationMemberParams{
		ActorID:        actorID,
		OrganizationID: orgID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
