package namespaces

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (o *Operations) Create(ctx context.Context, name string) (*models.Namespace, error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	project, err := models.NewNamespace(ident.Sub.ID, name)
	if err != nil {
		return nil, err
	}

	var created *models.Namespace
	err = database.RunTx(ctx, func(ctx context.Context) error {
		created, err = o.namespaces.Create(ctx, *project)
		if err != nil {
			return err
		}

		return o.namespaces.AddMember(ctx, models.NamespaceMember{
			UserID:      ident.Sub.ID,
			NamespaceID: created.ID,
			Role:        models.NamespaceMemberRoleOwner,
			AddedAt:     time.Now(),
			AddedBy:     ident.Sub.ID,
		})
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}
