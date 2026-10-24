package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.Project, error) {
	sqlcProjects, err := database.Queries(ctx, repo.q).ListProjectsByOrganizationID(ctx, &orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcProjects, mapProject), nil
}
