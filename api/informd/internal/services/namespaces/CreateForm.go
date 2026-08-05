package namespaces

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) CreateForm(ctx context.Context, title string, namespaceID uuid.UUID) (*models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.CreateForm")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := o.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckNamespace(ctx, ident.Sub.ID, namespace.ID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	form, err := models.NewForm(&namespaceID, namespace.OwnerID, ident.Sub.ID, title)
	if err != nil {
		return nil, err
	}

	created, err := o.forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}
