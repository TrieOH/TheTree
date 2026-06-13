package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMember, error) {
	ctx, span := repo.tracer.Start(ctx, "ListMembers")
	defer span.End()
	sqlcMembers, err := database.Queries(ctx, repo.q).ListOrganizationMembers(ctx, orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcMembers, mapOrganizationMember), nil
}
