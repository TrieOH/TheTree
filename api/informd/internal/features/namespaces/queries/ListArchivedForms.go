package queries

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/google/uuid"
)

func (s *Queries) ListArchivedForms(ctx context.Context, namespaceID uuid.UUID) (forms []models.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.ListArchivedForms")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var namespace *models.Namespace
	namespace, err = s.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	if sub.ID != namespace.OwnerID {
		_, err = s.namespaces.GetMember(ctx, sub.ID, namespaceID)
		if err != nil {
			return nil, err
		}
	}

	forms, err = s.forms.ListFromNamespaceArchived(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
