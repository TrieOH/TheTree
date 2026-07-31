package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListMembers")
	defer span.End()
	sqlcMembers, err := database.Queries(ctx, repo.q).ListOrganizationMembers(ctx, orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcMembers, mapOrganizationMember), nil
}
