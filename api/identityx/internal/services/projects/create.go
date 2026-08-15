package projects

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (o *Operations) Create(ctx context.Context, in models.CreateProjectInput) (*models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
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
	err = database.RunTx(ctx, func(ctx context.Context) error {
		created, err = o.projects.Create(ctx, *project)
		if err != nil {
			return err
		}

		// The Key-lifecycle module provisions the project's signing and
		// encryption keys inside the same tx; an org-created project (the
		// organizations feature) crosses the same seam, so every project
		// ships with keys regardless of which path created it.
		err = o.keys.Ensure(ctx, &created.ID)
		if err != nil {
			return err
		}

		member, err := models.NewProjectMember(created.ID, ident.Sub.ID, models.ProjectRoleOwner)
		if err != nil {
			return err
		}

		return o.projects.AddMember(ctx, *member)
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
