package projects

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (repo *Repo) Create(ctx context.Context, project models.Project) (*models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).CreateProject(ctx, sqlc.CreateProjectParams{
		OrganizationID: project.OrganizationID,
		OwnerID:        project.OwnerID,
		Name:           project.Name,
		Domain:         project.Domain,
		BrandSlug:      project.BrandSlug,
		Metadata:       project.Metadata,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProject(row)), nil
}
