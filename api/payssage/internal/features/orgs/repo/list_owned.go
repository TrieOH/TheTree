package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	sqlcOrgs, err := database.Queries(ctx, repo.q).ListOwnedOrganizations(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcOrgs, mapOrganization), nil
}
