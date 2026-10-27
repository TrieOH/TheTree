package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) CreateForm(ctx context.Context, title string, namespaceID uuid.UUID) (*models.Form, error) {
	namespace, err := o.Namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	ident, err := o.Authz.Namespace(ctx, namespace.ID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	form, err := models.NewForm(&namespaceID, namespace.OwnerID, ident.Sub.ID, title)
	if err != nil {
		return nil, err
	}

	created, err := o.Forms.Create(ctx, *form)
	if err != nil {
		return nil, err
	}

	return created, nil
}
