package commands

import (
	"Informd/models"
	"context"
	"lib/authz"
	"time"
)

func (s *Commands) Create(ctx context.Context, name string) (*models.Namespace, error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.Create")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	project, err := models.NewNamespace(sub.ID, name)
	if err != nil {
		return nil, err
	}

	var created *models.Namespace
	if err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = s.namespaces.Create(ctx, *project)
		if err != nil {
			return err
		}

		return s.namespaces.AddMember(ctx, models.NamespaceMember{
			UserID:      sub.ID,
			NamespaceID: created.ID,
			Role:        models.NamespaceMemberRoleOwner,
			AddedAt:     time.Now(),
			AddedBy:     sub.ID,
		})
	}); err != nil {
		return nil, err
	}

	return created, nil
}
