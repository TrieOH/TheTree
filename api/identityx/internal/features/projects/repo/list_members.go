package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error) {
	sqlcMembers, err := database.Queries(ctx, repo.q).ListProjectMembers(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcMembers, mapProjectMember), nil
}
