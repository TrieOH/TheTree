package repos

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (repo *Repo) AddMember(ctx context.Context, toCreate models.ProjectMember) error {
	ctx, span := telemetry.StartSpan(ctx, "AddMember")
	defer span.End()

	err := database.Queries(ctx, repo.q).AddProjectMember(ctx, sqlc.AddProjectMemberParams{
		ProjectID: toCreate.ProjectID,
		ActorID:   toCreate.ActorID,
		Role:      string(toCreate.Role),
		Metadata:  toCreate.Metadata,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
