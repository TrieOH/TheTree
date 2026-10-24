package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) ListForms(ctx context.Context, namespaceID uuid.UUID) (forms []models.Form, err error) {
	_, err = o.Authz.Namespace(ctx, namespaceID, models.NamespaceMemberRoleMember)
	if err != nil {
		return nil, err
	}

	forms, err = o.Forms.ListFromNamespace(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
