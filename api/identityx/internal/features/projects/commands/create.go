package commands

import (
	"IdentityX/models"
	"context"
	"time"
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
	if err = c.tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = c.projects.Create(ctx, *project)
		if err != nil {
			return err
		}

		member, err := models.NewProjectMember(created.ID, ident.Sub.ID, models.ProjectRoleOwner)
		if err != nil {
			return err
		}

		err = c.projects.AddMember(ctx, *member)
		if err != nil {
			return err
		}

		_, err = c.actors.Register(ctx, models.Actor{
			ProjectID:  &created.ID,
			AuthMethod: models.ApiKeyAuthMethod,
			VerifiedAt: new(time.Now()),
			Type:       models.ServiceActorType,
		})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return created, nil
}
