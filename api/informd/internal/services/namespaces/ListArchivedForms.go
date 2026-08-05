package namespaces

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) ListArchivedForms(ctx context.Context, namespaceID uuid.UUID) (forms []models.Form, err error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.ListArchivedForms")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckNamespace(ctx, ident.Sub.ID, namespaceID, models.NamespaceMemberRoleMember)
	if err != nil {
		return nil, err
	}

	forms, err = o.forms.ListFromNamespaceArchived(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
