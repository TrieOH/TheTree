package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	sqlcOrgs, err := database.Queries(ctx, repo.q).ListJoinedOrganizations(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcOrgs, mapOrganization), nil
}
