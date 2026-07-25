package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
	"lib/telemetry"
)

func (repo *Repo) ListMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListMembers")
	defer span.End()

	sqlcMembers, err := database.Queries(ctx, repo.q).ListProjectMembers(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcMembers, mapProjectMember), nil
}
