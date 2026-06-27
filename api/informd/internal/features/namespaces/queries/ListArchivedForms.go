package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/google/uuid"
)

func (q *Queries) ListArchivedForms(ctx context.Context, namespaceID uuid.UUID) (forms []models.Form, err error) {
	ctx, span := q.tracer.Start(ctx, "NamespaceService.ListArchivedForms")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var namespace *models.Namespace
	namespace, err = q.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != namespace.OwnerID {
		_, err = q.namespaces.GetMember(ctx, ident.Sub.ID, namespaceID)
		if err != nil {
			return nil, err
		}
	}

	forms, err = q.forms.ListFromNamespaceArchived(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
