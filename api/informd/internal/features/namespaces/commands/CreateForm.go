package commands

import (
	"context"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (s *Commands) CreateForm(ctx context.Context, title string, namespaceID uuid.UUID) (*models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.CreateForm")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := s.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckNamespace(ctx, ident.Sub.ID, namespace.ID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	form, err := models.NewForm(&namespaceID, namespace.OwnerID, ident.Sub.ID, title)
	if err != nil {
		return nil, err
	}

	created, err := s.forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}
