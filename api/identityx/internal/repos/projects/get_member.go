package projects

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetMember(ctx context.Context, actorID, projectID uuid.UUID) (*models.ProjectMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetMember")
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
