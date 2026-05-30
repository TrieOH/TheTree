package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error) {
	ctx, span := repo.tracer.Start(ctx, "ListMembers")
	defer span.End()
	sqlcMembers, err := database.Queries(ctx, repo.q).ListProjectMembers(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcMembers, mapProjectMember), nil
}
