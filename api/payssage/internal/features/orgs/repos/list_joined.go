package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListJoined")
	defer span.End()
	sqlcOrgs, err := database.Queries(ctx, repo.q).ListJoinedOrganizations(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcOrgs, mapOrganization), nil
}
