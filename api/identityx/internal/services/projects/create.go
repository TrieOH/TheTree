package projects

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
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

		var signKey *crypto.KeyPair
		signKey, err = crypto.GenerateKeyPair("signing")
		if err != nil {
			return err
		}

		_, err = o.keys.Create(ctx, &created.ID, signKey, "signing")
		if err != nil {
			return err
		}

		var encKey *crypto.KeyPair
		encKey, err = crypto.GenerateKeyPair("encryption")
		if err != nil {
			return err
		}

		_, err = o.keys.Create(ctx, &created.ID, encKey, "encryption")
		if err != nil {
			return err
		}

		var member *models.ProjectMember
		member, err = models.NewProjectMember(created.ID, ident.Sub.ID, models.ProjectRoleOwner)
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
