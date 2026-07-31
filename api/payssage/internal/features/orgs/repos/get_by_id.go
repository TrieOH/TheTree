package repos

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, orgID uuid.UUID) (*models.Organization, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByID")
	defer span.End()
	sqlcOrg, err := database.Queries(ctx, repo.q).GetOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOrganization(sqlcOrg)), nil
}
