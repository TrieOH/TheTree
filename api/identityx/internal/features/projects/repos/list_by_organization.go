package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ListByOrganization")
	defer span.End()
	sqlcProjects, err := database.Queries(ctx, repo.q).ListProjectsByOrganizationID(ctx, &orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcProjects, mapProject), nil
}
