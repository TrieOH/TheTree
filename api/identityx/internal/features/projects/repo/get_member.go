package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetMember(ctx context.Context, actorID, projectID uuid.UUID) (*models.ProjectMember, error) {
	sqlcMember, err := database.Queries(ctx, repo.q).GetProjectMemberByID(ctx, sqlc.GetProjectMemberByIDParams{
		ActorID:   actorID,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectMember(sqlcMember)), nil
}
