package repo

import (
	"IdentityX/internal/sqlc"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, actorID, projectID uuid.UUID) error {
	err := database.Queries(ctx, repo.q).RemoveProjectMember(ctx, sqlc.RemoveProjectMemberParams{
		ActorID:   actorID,
		ProjectID: projectID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
