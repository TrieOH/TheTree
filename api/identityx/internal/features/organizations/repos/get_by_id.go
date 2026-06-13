package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByID")
	defer span.End()
	sqlcOrg, err := database.Queries(ctx, repo.q).GetOrganizationByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOrganization(sqlcOrg)), nil
}
