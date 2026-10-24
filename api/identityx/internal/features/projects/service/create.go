package service

import (
	"IdentityX/models"
	"context"
)

func (o *Operations) Create(ctx context.Context, in models.CreateProjectInput) (*models.Project, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	project, err := models.NewProject(ident.Sub.ID, in.BrandSlug, in.Name, in.Domain, nil)
	if err != nil {
		return nil, err
	}

	var created *models.Project
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Projects.Create(ctx, *project)
		if err != nil {
			return err
		}

		// The Key-lifecycle module provisions the project's signing and
		// encryption keys inside the same tx; an org-created project (the
		// organizations feature) crosses the same seam, so every project
		// ships with keys regardless of which path created it.
		err = o.Keys.Ensure(ctx, &created.ID)
		if err != nil {
			return err
		}

		member, err := models.NewProjectMember(created.ID, ident.Sub.ID, models.ProjectRoleOwner)
		if err != nil {
			return err
		}

		return o.Projects.AddMember(ctx, *member)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
