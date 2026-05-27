package queries

import (
	models2 "IdentityX/models"
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

func (s *QueryService) ListForms(ctx context.Context, namespaceID uuid.UUID) (forms []models.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.ListForms")
	defer span.End()

	var sub *models2.UserSubject
	sub, err = models2.RequireSubject(ctx)
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

	forms, err = s.forms.ListFromNamespace(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
