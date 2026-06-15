package repos

import (
	"IdentityX/internal/database/sqlc"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) RemoveMember(ctx context.Context, actorID, projectID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RemoveMember")
	defer span.End()
	err := database.Queries(ctx, repo.q).RemoveProjectMember(ctx, sqlc.RemoveProjectMemberParams{
		ActorID:   actorID,
		ProjectID: projectID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
