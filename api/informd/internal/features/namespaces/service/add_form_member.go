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

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) AddFormMember(ctx context.Context, payload models.AddNamespaceFormMemberInput) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	if ident.Sub.ID == payload.UserID {
		return errs.ErrAddSelfToForm
	}

	namespace, err := o.Namespaces.GetByID(ctx, payload.NamespaceID)
	if err != nil {
		return err
	}

	if payload.UserID == namespace.OwnerID {
		return errs.ErrNamespaceOwnerAlreadyInForm
	}

	_, err = o.Authz.Namespace(ctx, namespace.ID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return err
	}

	_, err = o.Namespaces.GetMember(ctx, payload.UserID, namespace.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrNamespaceMemberAlreadyInForm
	}

	form, err := o.Forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	_, err = o.Forms.GetMember(ctx, payload.UserID, form.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if err == nil {
		return errs.ErrAlreadyFormMember
	}

	newMember := models.FormMember{
		UserID:  payload.UserID,
		FormID:  form.ID,
		Role:    payload.Role,
		AddedAt: time.Now(),
		AddedBy: ident.Sub.ID,
	}

	return o.Forms.AddMember(ctx, newMember)
}
