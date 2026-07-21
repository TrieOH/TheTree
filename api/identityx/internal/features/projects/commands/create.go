package commands

import (
	"IdentityX/models"
	"context"
)

func (c *Commands) Create(ctx context.Context, in models.CreateProjectInput) (*models.Project, error) {
	ctx, span := c.tracer.Start(ctx, "Create")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	project, err := models.NewProject(ident.Sub.ID, in.BrandSlug, in.Name, in.Domain, nil)
	if err != nil {
		return nil, err
	}

	var created *models.Project
	err = c.tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = c.projects.Create(ctx, *project)
		if err != nil {
			return err
		}

		member, err := models.NewProjectMember(created.ID, ident.Sub.ID, models.ProjectRoleOwner)
		if err != nil {
			return err
		}

		return c.projects.AddMember(ctx, *member)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
