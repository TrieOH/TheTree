package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	ctx, span := repo.tracer.Start(ctx, "ListJoined")
	defer span.End()
	sqlcOrgs, err := database.Queries(ctx, repo.q).ListJoinedOrganizations(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcOrgs, mapOrganization), nil
}
