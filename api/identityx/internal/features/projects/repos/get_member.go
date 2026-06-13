package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetMember(ctx context.Context, actorID, projectID uuid.UUID) (*models.ProjectMember, error) {
	ctx, span := repo.tracer.Start(ctx, "GetMember")
	defer span.End()
	sqlcMember, err := database.Queries(ctx, repo.q).GetProjectMemberByID(ctx, sqlc.GetProjectMemberByIDParams{
		ActorID:   actorID,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectMember(sqlcMember)), nil
}
