package service

import (
	"context"
	"errors"

	"Informd/errs"
	"Informd/models"
	"lib/database"
)

func (o *Operations) RemoveMember(ctx context.Context, payload models.RemoveNamespaceMemberInput) error {
	namespace, err := o.Namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Namespace(ctx, payload.NamespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Namespaces.GetMember(ctx, payload.UserID, namespace.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err != nil {
		return errs.ErrNotNamespaceMember
	}

	return o.Namespaces.RemoveMember(ctx, payload.UserID, payload.NamespaceID)
}
