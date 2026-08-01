package projects

import (
	"IdentityX/internal/sqlc"
	"context"
	"lib/database"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveMember(ctx context.Context, actorID, projectID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "RemoveMember")
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
