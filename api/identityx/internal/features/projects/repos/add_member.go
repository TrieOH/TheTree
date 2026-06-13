package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) AddMember(ctx context.Context, toCreate models.ProjectMember) error {
	ctx, span := repo.tracer.Start(ctx, "AddMember")
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
