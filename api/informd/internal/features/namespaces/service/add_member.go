package service

import (
	"context"
	"errors"
	idx "sdk/identityx"
	"time"

	"Informd/errs"
	"Informd/models"
	"lib/database"
)

func (o *Operations) AddMember(ctx context.Context, payload models.AddNamespaceMemberInput) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return errs.ErrAddSelfToNamespace
	}

	namespace, err := o.Namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return errs.ErrAddOwnerToNamespace
	}

	_, err = o.Authz.Namespace(ctx, payload.NamespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Namespaces.GetMember(ctx, payload.UserID, namespace.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyNamespaceMember
	}

	newMember := models.NamespaceMember{
		UserID:      payload.UserID,
		NamespaceID: payload.NamespaceID,
		Role:        payload.Role,
		AddedAt:     time.Now(),
		AddedBy:     ident.Sub.ID,
	}

	err = o.Namespaces.AddMember(ctx, newMember)
	if err != nil {
		return err
	}
	return nil
}
