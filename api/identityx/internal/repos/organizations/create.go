package organizations

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Organization) (*models.Organization, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()
	sqlcOrg, err := database.Queries(ctx, repo.q).CreateOrganization(ctx, sqlc.CreateOrganizationParams{
		OwnerID:  toCreate.OwnerID,
		Name:     toCreate.Name,
		Slug:     toCreate.Slug,
		Metadata: toCreate.Metadata,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOrganization(sqlcOrg)), nil
}
