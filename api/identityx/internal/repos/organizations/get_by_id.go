package organizations

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()
	sqlcOrg, err := database.Queries(ctx, repo.q).GetOrganizationByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOrganization(sqlcOrg)), nil
}
