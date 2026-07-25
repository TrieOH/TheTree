package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
	"lib/telemetry"
)

func (repo *Repo) ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByOrganization")
	defer span.End()

	sqlcProjects, err := database.Queries(ctx, repo.q).ListProjectsByOrganizationID(ctx, &orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcProjects, mapProject), nil
}
