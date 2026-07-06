package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, project models.Project) (*models.Project, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Create")
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
