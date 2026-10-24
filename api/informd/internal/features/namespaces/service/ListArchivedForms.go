package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) ListArchivedForms(ctx context.Context, namespaceID uuid.UUID) (forms []models.Form, err error) {
	_, err = o.Authz.Namespace(ctx, namespaceID, models.NamespaceMemberRoleMember)
	if err != nil {
		return nil, err
	}

	forms, err = o.Forms.ListFromNamespaceArchived(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
